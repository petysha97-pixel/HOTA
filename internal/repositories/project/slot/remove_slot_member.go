package slot

import (
	"HOTA/internal/models"
	"errors"
	"fmt"
)

// Снимаем утверждённого разработчика и открываем слот.
func RemoveApprovedMember(slotID int) (int64, error) {
	// Начинаем общую операцию из нескольких SQL-запросов.
	tx, err := models.UserDB.Begin()
	if err != nil {
		return 0, fmt.Errorf("начало транзакции: %w", err)
	}
	defer tx.Rollback()

	// Снимаем утверждённую заявку.
	res, err := tx.Exec(`
		UPDATE applicationsSlot
		SET status = 'removed',
		    updated_at = CURRENT_TIMESTAMP
		WHERE slot_id = ?
		  AND status = 'approved'
	`, slotID)
	if err != nil {
		return 0, fmt.Errorf("снятие заявки: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("подсчёт изменённых заявок: %w", err)
	}

	if count == 0 {
		return 0, errors.New("нет утверждённого разработчика для снятия")
	}

	if count != 1 {
		return 0, errors.New("на слоте несколько утверждённых заявок")
	}

	// Убираем исполнителя и открываем слот.
	res, err = tx.Exec(`
		UPDATE slots
		SET user_id = NULL,
		    status = 'open'
		WHERE id = ?
		  AND status IN ('open', 'review')
	`, slotID)
	if err != nil {
		return 0, fmt.Errorf("освобождение слота: %w", err)
	}

	changed, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("подсчёт изменённых слотов: %w", err)
	}

	if changed != 1 {
		return 0, errors.New("слот не найден или его нельзя освободить")
	}

	// Только здесь сохраняются оба изменения.
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("сохранение транзакции: %w", err)
	}

	return count, nil
}
