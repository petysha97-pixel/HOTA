package slot

import (
	"HOTA/internal/models"
	"errors"
	"fmt"
)

var ErrInviteConflict = errors.New("не удалось пригласить: слот занят, разработчик недоступен или уже есть приглашение либо отклик")

// Приглашаем разработчика на свободный слот.
func InviteToSlot(slotID, applicantID, ownerID int) error {
	// Проверяем слот и сохраняем приглашение одним запросом.
	result, err := models.UserDB.Exec(`
  INSERT INTO invites (slot_id, user_id, status)
  SELECT s.id, ?, 'pending'
  FROM slots s
  JOIN projects p ON p.id = s.project_id
  WHERE s.id = ?
    AND s.status = 'open'
    AND s.user_id IS NULL
    AND p.owner_id = ?
    AND p.status != 'finished'
    AND p.owner_id != ?
    AND EXISTS (SELECT 1 FROM users WHERE id = ?)
    AND NOT EXISTS (
        SELECT 1 FROM invites i
        WHERE i.slot_id = s.id AND i.user_id = ?
          AND i.status = 'pending'
    )
    AND NOT EXISTS (
        SELECT 1 FROM applicationsSlot a
        WHERE a.slot_id = s.id AND a.user_id = ?
          AND a.status IN ('pending', 'approved')
    )
    AND NOT EXISTS (
        SELECT 1 FROM applicationsSlot a
        WHERE a.slot_id = s.id AND a.status = 'approved'
    )
 `, applicantID, slotID, ownerID, applicantID, applicantID, applicantID, applicantID)
	if err != nil {
		return fmt.Errorf("сохранение приглашения: %w", err)
	}

	// Проверяем, что приглашение добавилось.
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("проверка приглашения: %w", err)
	}
	if count == 0 {
		return ErrInviteConflict
	}
	return nil
}
