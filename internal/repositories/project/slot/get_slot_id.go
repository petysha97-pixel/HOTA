package slot

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	"database/sql"
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

// обновляем статус слота
func UpdateSlotByID(slotID int, status string) (*models.Slot, error) {
	query := `UPDATE slots SET status = ? WHERE id = ?`

	if _, err := models.UserDB.Exec(query, status, slotID); err != nil {
		return nil, fmt.Errorf("ошибка обновления слота: %w", err)
	}

	return GetSlotByID(slotID)
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
