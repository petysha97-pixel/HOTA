package slot

import (
	"HOTA/internal/models"
	"fmt"
)

// создание заявки
func CreateApplication(slotID, userID int) error {
	query := `INSERT INTO applications (slot_id, user_id, status) VALUES (?, ?, 'pending')`
	_, err := models.UserDB.Exec(query, slotID, userID)
	if err != nil {
		return fmt.Errorf("ошибка создания заявки: %w", err)
	}
	return nil
}

// проверка заявки на слот
func GetUserCheckingApplications(slotID, userID int) (*models.AplicationSlot, error) {

	query := `SELECT id, slot_id, user_id, status, creat_add FROM apllicationsSlot 
 WHERE slot_id = ?
 AND user_id = ?`

	var app models.AplicationSlot
	err := models.UserDB.QueryRow(query, slotID, userID).Scan(&app.ID, &app.SlotID, &app.UserID, &app.Status, &app.Creat_add)
	if err != nil {
		return nil, fmt.Errorf("заявка не найдена %w", err)
	}

	return &app, nil

}

// достаём все заявки по ID слота
func GetApplicationsBySlot(slotID int) ([]models.AplicationSlot, error) {
	query := `SELECT id, slot_id, user_id, status, creat_add FROM applications WHERE slot_id = ?`
	rows, err := models.UserDB.Query(query, slotID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заявок: %w", err)
	}
	defer rows.Close()

	var app []models.AplicationSlot
	for rows.Next() {
		var a models.AplicationSlot
		err := rows.Scan(&a.ID, &a.SlotID, &a.UserID, &a.Status, &a.Creat_add)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования: %w", err)
		}
		app = append(app, a)
	}
	return app, nil
}

// обновляем статус заявки
func UpdateApplicationStatus(appID int, status string) error {
	query := `UPDATE apllicationsSlot SET status = ? WHERE id = ?`
	_, err := models.UserDB.Exec(query, status, appID)
	if err != nil {
		return fmt.Errorf("ошибка обновления заявки: %w", err)
	}
	return nil
}

// отклоняем все остальные заявки на слот, кроме указанной
func RejectOtherApplications(slotID, excludeAppID int) error {
	query := `UPDATE apllicationsSlot SET status = 'rejected' WHERE slot_id = ? AND id != ? AND status = 'pending'`
	_, err := models.UserDB.Exec(query, slotID, excludeAppID)
	if err != nil {
		return fmt.Errorf("ошибка отклонения заявок: %w", err)
	}
	return nil
}
