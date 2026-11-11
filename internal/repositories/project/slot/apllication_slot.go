package slot

import (
	"HOTA/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

var ErrApplicationConflict = errors.New("не удалось подать отклик: слот или заявка изменились")

// Создаём отклик, если набор ещё открыт.
func CreateApplication(slotID, userID int, message string) error {
	// SELECT вернёт строку для вставки только при открытом наборе.
	result, err := models.UserDB.Exec(`
		INSERT INTO applicationsSlot (slot_id, user_id, status, message)
		SELECT s.id, ?, 'pending', ?
		FROM slots s
		JOIN projects p ON p.id = s.project_id
		WHERE s.id = ?
		  AND s.status IN ('open', 'review')
		  AND p.status != 'finished'
		ON CONFLICT(slot_id, user_id) DO NOTHING
	`, userID, message, slotID)
	if err != nil {
		return fmt.Errorf("создание заявки: %w", err)
	}

	// Если строка не добавилась, набор закрыт или отклик уже существует.
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("проверка создания заявки: %w", err)
	}
	if count == 0 {
		return ErrApplicationConflict
	}
	return nil
}

// Получаем ожидающую заявку кандидата.
func GetUserCheckingApplications(slotID, userID int) (*models.AplicationSlot, error) {
	query := `SELECT id, slot_id, user_id, status, message, created_at FROM applicationsSlot
	WHERE slot_id = ?
	AND user_id = ?
	AND status = 'pending'`

	var app models.AplicationSlot
	err := models.UserDB.QueryRow(query, slotID, userID).Scan(&app.ID, &app.SlotID, &app.UserID, &app.Status, &app.Message, &app.CreatAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска заявки: %w", err)
	}

	return &app, nil
}

// Получаем все заявки на слот.
func GetApplicationsBySlot(slotID int) ([]models.AplicationSlot, error) {
	// Синий значок получает кандидат, который принял приглашение на этот слот.
	query := `SELECT a.id, a.slot_id, a.user_id, a.status, a.message, a.created_at,
	EXISTS(SELECT 1 FROM invites i WHERE i.slot_id = a.slot_id AND i.user_id = a.user_id AND i.status = 'accepted')
	FROM applicationsSlot a WHERE a.slot_id = ?`
	rows, err := models.UserDB.Query(query, slotID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заявок: %w", err)
	}
	defer rows.Close()

	app := []models.AplicationSlot{}
	for rows.Next() {
		var a models.AplicationSlot
		err := rows.Scan(&a.ID, &a.SlotID, &a.UserID, &a.Status, &a.Message, &a.CreatAt, &a.Invited)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования: %w", err)
		}
		app = append(app, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения заявок: %w", err)
	}
	return app, nil
}

// Получаем прежнюю заявку разработчика на этот слот.
func GetUserApplication(slotID, userID int) (*models.AplicationSlot, error) {
	query := `SELECT id, slot_id, user_id, status, message, created_at
 FROM applicationsSlot WHERE slot_id = ? AND user_id = ?`

	var app models.AplicationSlot
	err := models.UserDB.QueryRow(query, slotID, userID).Scan(
		&app.ID, &app.SlotID, &app.UserID, &app.Status, &app.Message, &app.CreatAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска заявки: %w", err)
	}
	return &app, nil
}

// Возвращаем отклонённую или снятую заявку в очередь.
func ReopenApplication(appID int, message string) error {
	// Проверяем состояние заявки и слота в момент обновления.
	result, err := models.UserDB.Exec(`
		UPDATE applicationsSlot
		SET status = 'pending', message = ?,
		    created_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND status IN ('rejected', 'removed')
		  AND EXISTS (
		      SELECT 1
		      FROM slots s
		      JOIN projects p ON p.id = s.project_id
		      WHERE s.id = applicationsSlot.slot_id
		        AND s.status IN ('open', 'review')
		        AND p.status != 'finished'
		  )
	`, message, appID)
	if err != nil {
		return fmt.Errorf("повторная подача заявки: %w", err)
	}

	// Если условие перестало выполняться, отклик остаётся как был.
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("проверка обновления заявки: %w", err)
	}
	if count == 0 {
		return ErrApplicationConflict
	}
	return nil
}
