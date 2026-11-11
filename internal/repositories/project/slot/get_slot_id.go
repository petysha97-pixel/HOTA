package slot

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	"database/sql"
	"errors"
	"fmt"
)

// Получение слота по айди
func GetSlotByID(slotID int) (*models.Slot, error) {
	query := `SELECT id, project_id, name, description, rolle, user_id, status, created_at FROM slots WHERE id = ?`

	var slot models.Slot
	var userID sql.NullInt64 // исполнителя может не быть (NULL)

	err := models.UserDB.QueryRow(query, slotID).Scan(&slot.ID, &slot.ProjectID, &slot.Name, &slot.Description, &slot.Rolle, &userID, &slot.Status, &slot.CreatAt)
	if err != nil {
		return nil, fmt.Errorf("слот не найден %w", err)
	}
	slot.UserID = projectRepo.UserIDOrNil(userID)

	// стек слота лежит в отдельной таблице slot_stacks
	slot.StackID, err = projectRepo.GetStackIDsBySlotID(slot.ID)
	if err != nil {
		return nil, err
	}

	return &slot, nil
}

var ErrSlotStatusConflict = errors.New("состояние слота уже изменилось")

// Меняем статус слота. При переоткрытии снимаем исполнителя.
func UpdateSlotByID(slotID, ownerID int, oldStatus, newStatus string) error {
	// Сохраняем изменения слота и заявки вместе.
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("начало транзакции: %w", err)
	}
	defer tx.Rollback()

	// Меняем слот, только если его статус ещё тот, который прочитал обработчик.
	// При переходе в open очищаем исполнителя.
	result, err := tx.Exec(`
		UPDATE slots
		SET status = ?,
		    user_id = CASE WHEN ? = 'open' THEN NULL ELSE user_id END
		WHERE id = ?
		  AND status = ?
		  AND EXISTS (
		      SELECT 1 FROM projects p
		      WHERE p.id = slots.project_id
		        AND p.owner_id = ?
		        AND p.status != 'finished'
		  )
		  AND (
		      ? = 'open'
		      OR EXISTS (
		          SELECT 1 FROM applicationsSlot a
		          WHERE a.slot_id = slots.id
		            AND a.user_id = slots.user_id
		            AND a.status = 'approved'
		      )
		  )
	`, newStatus, newStatus, slotID, oldStatus, ownerID, newStatus)
	if err != nil {
		return fmt.Errorf("смена статуса слота: %w", err)
	}

	// Ноль строк означает, что переход уже нельзя выполнить.
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("проверка обновления слота: %w", err)
	}
	if count == 0 {
		return ErrSlotStatusConflict
	}

	// При переоткрытии помечаем принятую заявку как снятую.
	if oldStatus == "close" && newStatus == "open" {
		_, err = tx.Exec(`
			UPDATE applicationsSlot
			SET status = 'removed', updated_at = CURRENT_TIMESTAMP
			WHERE slot_id = ? AND status = 'approved'
		`, slotID)
		if err != nil {
			return fmt.Errorf("снятие исполнителя: %w", err)
		}
	}

	// Фиксируем смену статуса и снятие исполнителя.
	return tx.Commit()
}

// обновляем название, описание, роль и стек слота
func UpdateSlotInfo(slotID int, name, description, rolle string, stackIDs []int) error {
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() // откатится, если не закоммитим

	_, err = tx.Exec(`UPDATE slots SET name = ?, description = ?, rolle = ? WHERE id = ?`, name, description, rolle, slotID)
	if err != nil {
		return fmt.Errorf("ошибка обновления слота: %w", err)
	}

	// стек слота перезаписываем целиком
	if _, err = tx.Exec(`DELETE FROM slot_stacks WHERE slot_id = ?`, slotID); err != nil {
		return fmt.Errorf("ошибка удаления стеков слота: %w", err)
	}
	for _, stackID := range stackIDs {
		if _, err = tx.Exec(`INSERT INTO slot_stacks (slot_id, stack_id) VALUES (?, ?)`, slotID, stackID); err != nil {
			return fmt.Errorf("ошибка добавления стека в слот: %w", err)
		}
	}

	return tx.Commit()
}

// удаляем слот вместе с его стеком и заявками
func DeleteSlotByID(slotID int) error {
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() // откатится, если не закоммитим

	if _, err = tx.Exec(`DELETE FROM applicationsSlot WHERE slot_id = ?`, slotID); err != nil {
		return fmt.Errorf("ошибка удаления заявок слота: %w", err)
	}
	if _, err = tx.Exec(`DELETE FROM slot_stacks WHERE slot_id = ?`, slotID); err != nil {
		return fmt.Errorf("ошибка удаления стеков слота: %w", err)
	}

	res, err := tx.Exec(`DELETE FROM slots WHERE id = ?`, slotID)
	if err != nil {
		return fmt.Errorf("ошибка удаления слота: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("слот %d не найден", slotID)
	}

	return tx.Commit()
}
