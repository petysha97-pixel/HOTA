package slot

import (
	"HOTA/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

// создание заявки
func CreateApplication(slotID, userID int) error {
	query := `INSERT INTO applicationsSlot (slot_id, user_id, status) VALUES (?, ?, 'pending')`
	_, err := models.UserDB.Exec(query, slotID, userID)
	if err != nil {
		return fmt.Errorf("ошибка создания заявки: %w", err)
	}
	return nil
}

// ищем заявку пользователя на слот, которая ещё на рассмотрении
func GetUserCheckingApplications(slotID, userID int) (*models.AplicationSlot, error) {
	query := `SELECT id, slot_id, user_id, status, created_at FROM applicationsSlot
	WHERE slot_id = ?
	AND user_id = ?
	AND status = 'pending'`

	var app models.AplicationSlot
	err := models.UserDB.QueryRow(query, slotID, userID).Scan(&app.ID, &app.SlotID, &app.UserID, &app.Status, &app.Creat_add)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска заявки: %w", err)
	}

	return &app, nil
}

// достаём все заявки по ID слота
func GetApplicationsBySlot(slotID int) ([]models.AplicationSlot, error) {
	query := `SELECT id, slot_id, user_id, status, created_at FROM applicationsSlot WHERE slot_id = ?`
	rows, err := models.UserDB.Query(query, slotID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заявок: %w", err)
	}
	defer rows.Close()

	app := []models.AplicationSlot{}
	for rows.Next() {
		var a models.AplicationSlot
		err := rows.Scan(&a.ID, &a.SlotID, &a.UserID, &a.Status, &a.Creat_add)
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

// обновляем статус заявки
func UpdateApplicationStatus(appID int, status string) error {
	query := `UPDATE applicationsSlot SET status = ? WHERE id = ?`
	_, err := models.UserDB.Exec(query, status, appID)
	if err != nil {
		return fmt.Errorf("ошибка обновления заявки: %w", err)
	}
	return nil
}

// отклоняем все остальные заявки на слот, кроме указанной
func RejectOtherApplications(slotID, excludeAppID int) error {
	query := `UPDATE applicationsSlot SET status = 'rejected' WHERE slot_id = ? AND id != ? AND status = 'pending'`
	_, err := models.UserDB.Exec(query, slotID, excludeAppID)
	if err != nil {
		return fmt.Errorf("ошибка отклонения заявок: %w", err)
	}
	return nil
}

// проверка на то, что в слоте утверждён разработчик
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

// при переоткрытии снимаем принятую заявку
func ResetApprovedApplication(slotID int) error {
	query := `UPDATE applicationsSlot SET status = 'rejected' WHERE slot_id = ? AND status = 'approved'`
	_, err := models.UserDB.Exec(query, slotID)
	if err != nil {
		return fmt.Errorf("ошибка снятия заявки: %w", err)
	}
	return nil
}

// снимаем утвержденного разработчика со слота
func RemoveApprovedMember(slotID int) (int64, error) {
	query := `UPDATE applicationsSlot SET status = 'removed' WHERE slot_id = ? AND status = 'approved'`

	res, err := models.UserDB.Exec(query, slotID)
	if err != nil {
		return 0, fmt.Errorf("ошибка снятия разработчика %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("ошибка подсчета снятия разработчия со слота %w", err)
	}
	if count == 0 {
		return 0, errors.New("нет утверждённого разработчика для снятия")
	}

	return count, nil

}
