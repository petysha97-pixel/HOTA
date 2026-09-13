package slot

import (
	"HOTA/internal/models"
	"fmt"
)

// Получение слота по айди
func GetSlotByID(slotID int) (*models.Slot, error) {
	query := `SELECT id, project_id, rolle, status, created_at FROM slots WHERE id = ?`
	var slot models.Slot
	err := models.UserDB.QueryRow(query, slotID).Scan(&slot.ID, &slot.ProjectID, &slot.Rolle, &slot.Status, &slot.CreatAt)
	if err != nil {
		return nil, fmt.Errorf("слот не найден %w", err)
	}

	return &slot, nil
}

// обновляем статус слота
func UpdateSlotByID(slotID int, status string) (*models.Slot, error) {
	query := `UPDATE slots SET status = ? WHERE id = ? RETURNING id, project_id, rolle, status, created_at`
    

	var slot models.Slot
	err := models.UserDB.QueryRow(query, status, slotID).Scan(
		&slot.ID, &slot.ProjectID, &slot.Rolle, &slot.Status, &slot.CreatAt,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления слота: %w", err)
	}

	return &slot, nil
}
