package slot

import (
	"HOTA/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

var ErrInviteNotFound = errors.New("приглашение не найдено")
var ErrInviteForbidden = errors.New("это приглашение другому разработчику")
var ErrInviteAnswered = errors.New("на приглашение уже ответили")
var ErrInviteUnavailable = errors.New("набор на слот уже закрыт")

// Сохраняем ответ кандидата.
// После принятия добавляем его в очередь, знакомство начинает владелец.
func AnswerInvite(inviteID, userID int, accept bool) error {
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("начало транзакции: %w", err)
	}
	defer tx.Rollback()

	status := "declined"
	if accept {
		status = "accepted"
	}

	// Сначала меняем приглашение. Это также блокирует конкурирующие ответы в SQLite.
	// Если дальше возникнет ошибка, ответ откатится вместе с заявкой.
	result, err := tx.Exec(`
  UPDATE invites SET status = ?
  WHERE id = ? AND user_id = ? AND status = 'pending'
 `, status, inviteID, userID)
	if err != nil {
		return fmt.Errorf("ответ на приглашение: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("проверка ответа: %w", err)
	}
	if count == 0 {
		// Различаем чужое приглашение, повторный ответ и неверный ID.
		var recipientID int
		err = tx.QueryRow("SELECT user_id FROM invites WHERE id = ?", inviteID).Scan(&recipientID)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInviteNotFound
		}
		if err != nil {
			return fmt.Errorf("проверка приглашения: %w", err)
		}
		if recipientID != userID {
			return ErrInviteForbidden
		}
		return ErrInviteAnswered
	}

	if !accept {
		// При отказе заявка в очередь не добавляется.
		return tx.Commit()
	}

	// Проверяем, что за время ожидания слот не закрылся.
	var slotID int
	var available bool
	err = tx.QueryRow(`
  SELECT i.slot_id,
   s.status IN ('open', 'review') AND p.status != 'finished'
   AND NOT EXISTS (
       SELECT 1 FROM applicationsSlot a
       WHERE a.slot_id = s.id AND a.status = 'approved'
   )
  FROM invites i
  JOIN slots s ON s.id = i.slot_id
  JOIN projects p ON p.id = s.project_id
  WHERE i.id = ?
 `, inviteID).Scan(&slotID, &available)
	if err != nil {
		return fmt.Errorf("проверка слота приглашения: %w", err)
	}
	if !available {
		return ErrInviteUnavailable
	}

	// Разработчик мог откликнуться сам, пока приглашение ожидало ответа.
	var previousStatus string
	err = tx.QueryRow(`
  SELECT status FROM applicationsSlot WHERE slot_id = ? AND user_id = ?
 `, slotID, userID).Scan(&previousStatus)
	if errors.Is(err, sql.ErrNoRows) {
		// После согласия создаём обычную ожидающую заявку.
		_, err = tx.Exec(`
   INSERT INTO applicationsSlot(slot_id, user_id, status, message)
   VALUES(?, ?, 'pending', '')
  `, slotID, userID)
	} else if err == nil {
		switch previousStatus {
		case "pending":
			// Уже есть отклик в очереди, сохраняем его сообщение и дату.
		case "rejected", "removed":
			// Возвращаем прежнюю заявку в очередь после согласия кандидата.
			_, err = tx.Exec(`
    UPDATE applicationsSlot
    SET status = 'pending', created_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
    WHERE slot_id = ? AND user_id = ?
   `, slotID, userID)
		default:
			return ErrInviteUnavailable
		}
	}
	if err != nil {
		return fmt.Errorf("добавление кандидата в очередь: %w", err)
	}

	// Ответ и заявка сохраняются вместе. Слот остаётся в прежнем статусе.
	return tx.Commit()
}
