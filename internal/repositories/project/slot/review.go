package slot

import (
	"HOTA/internal/models"
	"errors"
	"fmt"
)

var ErrReviewConflict = errors.New("слот занят или нет ожидающей заявки")

// Берём кандидата на знакомство и уведомляем его.
func ReviewSlotApplication(slotID, applicantID, ownerID int) error {
	// Сохраняем знакомство и уведомление вместе.
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("начало транзакции: %w", err)
	}
	defer tx.Rollback()

	// Назначаем кандидата, если слот свободен и его заявка ещё ожидает решения.
	result, err := tx.Exec(`
		UPDATE slots
		SET user_id = ?, status = 'review'
		WHERE id = ?
		  AND status = 'open'
		  AND user_id IS NULL
		  AND EXISTS (
		      SELECT 1 FROM projects p
		      WHERE p.id = slots.project_id
		        AND p.owner_id = ?
		        AND p.status != 'finished'
		  )
		  AND EXISTS (
		      SELECT 1 FROM applicationsSlot a
		      WHERE a.slot_id = slots.id
		        AND a.user_id = ?
		        AND a.status = 'pending'
		  )
		  AND NOT EXISTS (
		      SELECT 1 FROM applicationsSlot a
		      WHERE a.slot_id = slots.id
		        AND a.status = 'approved'
		  )
	`, applicantID, slotID, ownerID, applicantID)
	if err != nil {
		return fmt.Errorf("назначение кандидата на знакомство: %w", err)
	}

	// Если ничего не изменилось, условия для знакомства уже не выполняются.
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("проверка назначения кандидата: %w", err)
	}
	if count == 0 {
		return ErrReviewConflict
	}

	// Сообщаем выбранному разработчику, на какой слот его взяли.
	_, err = tx.Exec(`
		INSERT INTO notifications(user_id, text, project_id)
		SELECT s.user_id,
		    'Вас взяли на знакомство в проект «' || p.name || '». Слот «' || s.name || '».',
		    p.id
		FROM slots s
		JOIN projects p ON p.id = s.project_id
		WHERE s.id = ?
	`, slotID)
	if err != nil {
		return fmt.Errorf("уведомление о знакомстве: %w", err)
	}

	return tx.Commit()
}
