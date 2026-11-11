package slot

import (
	"HOTA/internal/models"
	"errors"
	"fmt"
)

var ErrApprovalConflict = errors.New("слот занят или заявка уже изменилась")

// Утверждаем кандидата и закрываем слот.
// Остальные заявки в очереди отклоняем.
func ApproveSlotApplication(slotID, applicantID, ownerID int) error {
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("начало транзакции: %w", err)
	}
	defer tx.Rollback()

	// Утверждаем только текущего кандидата на знакомстве.
	// Повторно проверяем владельца и состояние проекта.
	res, err := tx.Exec(`
	UPDATE slots
	SET user_id = ?,
    status = 'close'
	WHERE id = ?
	  AND status = 'review' AND user_id = ?
	  AND EXISTS (
	      SELECT 1 FROM projects
	      WHERE id = slots.project_id
	        AND owner_id = ?
	        AND status != 'finished'
	  )
`, applicantID, slotID, applicantID, ownerID)
	if err != nil {
		return fmt.Errorf("назначение исполнителя: %w", err)
	}

	changed, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("подсчёт изменённых слотов: %w", err)
	}
	if changed != 1 {
		return ErrApprovalConflict
	}

	// Утверждаем только заявку, которая ещё ждёт решения.
	res, err = tx.Exec(`
		UPDATE applicationsSlot
		SET status = 'approved', updated_at = CURRENT_TIMESTAMP
		WHERE slot_id = ? AND user_id = ? AND status = 'pending'
	`, slotID, applicantID)
	if err != nil {
		return fmt.Errorf("утверждение заявки: %w", err)
	}

	changed, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("подсчёт изменённых заявок: %w", err)
	}
	if changed != 1 {
		return ErrApprovalConflict
	}

	// Отклоняем оставшиеся заявки, которые ещё ждут решения.
	_, err = tx.Exec(`
		UPDATE applicationsSlot
		SET status = 'rejected', updated_at = CURRENT_TIMESTAMP
		WHERE slot_id = ? AND status = 'pending'
	`, slotID)
	if err != nil {
		return fmt.Errorf("отклонение остальных заявок: %w", err)
	}

	// Закрываем оставшиеся приглашения, чтобы они не ожили при переоткрытии слота.
	_, err = tx.Exec(`UPDATE invites SET status = 'cancelled' WHERE slot_id = ? AND status = 'pending'`, slotID)
	if err != nil {
		return fmt.Errorf("отмена оставшихся приглашений: %w", err)
	}

	return tx.Commit()
}

// Проверяем, есть ли утверждённый разработчик на слоте.
func ApprovedAplecation(slotID int) (bool, error) {
	query := `SELECT COUNT(*) FROM applicationsSlot
	WHERE slot_id = ?
	AND status = 'approved'`

	var count int

	err := models.UserDB.QueryRow(query, slotID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("ошибка проверки утверждённого разработчика на слот: %w", err)
	}

	return count > 0, nil
}

// Отклоняем остальные ожидающие заявки на слот.
func RejectOtherApplications(slotID, excludeAppID int) error {
	query := `UPDATE applicationsSlot SET status = 'rejected', updated_at = CURRENT_TIMESTAMP WHERE slot_id = ? AND id != ? AND status = 'pending'`
	_, err := models.UserDB.Exec(query, slotID, excludeAppID)
	if err != nil {
		return fmt.Errorf("ошибка отклонения заявок: %w", err)
	}
	return nil
}
