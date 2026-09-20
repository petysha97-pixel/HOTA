package project

import (
	"HOTA/internal/models"
	"fmt"
)

// обновить статус проекта
func UpdateProjectStatus(roomID int, status string) error {
	query := `UPDATE projects SET status = ? WHERE id = ?`
	_, err := models.UserDB.Exec(query, status, roomID)
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса: %w", err)
	}
	return nil
}

// считаем количество слотов в стататусах в проекте
func CountSlotStatus(IDproject int) (totalSlot, openSlot, closeSlot, doneSlot int, err error) {
	
	var (
	total int
	open int
	close int
	done int
	)

	query := `SELECT COUNT(*) AS total_slots, 
	SUM(CASE WHEN status = 'open' THEN 1 ELSE 0 END) AS open_slots,
	SUM(CASE WHEN status = 'close' THEN 1 ELSE 0 END) AS close_slots,
	SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END) AS done_slots,
	FROM slots
	WHERE 
    project_id = ?;
	`

	err = models.UserDB.QueryRow(query, IDproject).Scan(&total, &open, &close, &done)
if err != nil {
	return 0, 0, 0, 0, fmt.Errorf("ошибка подчсета слотов %w", err)
}

return total, open, close, done, nil

}
