package project

import (
	"HOTA/internal/models"
	"fmt"
)

// обновить статус комнаты
func UpdateProjectStatus(roomID int, status string) error {
	query := `UPDATE projects SET status = ? WHERE id = ?`
	_, err := models.UserDB.Exec(query, status, roomID)
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса: %w", err)
	}
	return nil
}
