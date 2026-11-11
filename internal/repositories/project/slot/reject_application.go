package slot

import (
	"HOTA/internal/models"
	"errors"
	"fmt"
)

var ErrRejectConflict = errors.New("заявка или слот уже изменились")

// Отклоняем заявку.
// Если кандидат был на знакомстве, освобождаем слот.
func RejectSlotApplication(slotID, applicantID, ownerID int) error {
	// Сохраняем отказ и освобождение слота вместе.
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("начало транзакции: %w", err)
	}
	defer tx.Rollback()

	// Отклоняем только ожидающую заявку в доступном владельцу слоте.
	result, err := tx.Exec(`
		UPDATE applicationsSlot
		SET status = 'rejected', updated_at = CURRENT_TIMESTAMP
		WHERE slot_id = ?
		  AND user_id = ?
		  AND status = 'pending'
		  AND EXISTS (
		      SELECT 1
		      FROM slots s
		      JOIN projects p ON p.id = s.project_id
		      WHERE s.id = applicationsSlot.slot_id
		        AND s.status IN ('open', 'review')
		        AND p.owner_id = ?
		        AND p.status != 'finished'
		  )
	`, slotID, applicantID, ownerID)
	if err != nil {
		return fmt.Errorf("отклонение заявки: %w", err)
	}

	// Проверяем, что заявка действительно отклонена.
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("проверка отказа: %w", err)
	}
	if count != 1 {
		return ErrRejectConflict
	}

	// Освобождаем слот, только если отказали кандидату на знакомстве.
	_, err = tx.Exec(`
		UPDATE slots
		SET status = 'open', user_id = NULL
		WHERE id = ?
		  AND status = 'review'
		  AND user_id = ?
	`, slotID, applicantID)
	if err != nil {
		return fmt.Errorf("освобождение слота: %w", err)
	}

	// Фиксируем отказ и изменение слота.
	return tx.Commit()
}

// Меняем статус заявки.
func UpdateApplicationStatus(appID int, status string) error {
	query := `UPDATE applicationsSlot SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := models.UserDB.Exec(query, status, appID)
	if err != nil {
		return fmt.Errorf("ошибка обновления заявки: %w", err)
	}
	return nil
}
