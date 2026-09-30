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

	
	query := `SELECT
	COUNT(*) AS total_slots,
	COALESCE(SUM(CASE WHEN status = 'open' THEN 1 ELSE 0 END), 0) AS open_slots,
	COALESCE(SUM(CASE WHEN status = 'close' THEN 1 ELSE 0 END), 0) AS close_slots,
	COALESCE(SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END), 0) AS done_slots
	FROM slots
	WHERE project_id = ?`

	err = models.UserDB.QueryRow(query, IDproject).Scan(&totalSlot, &openSlot, &closeSlot, &doneSlot)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("ошибка подсчёта слотов: %w", err)
	}

	return totalSlot, openSlot, closeSlot, doneSlot, nil
}

// обновляем статус приватности проекта
func UpdateProjectPrivacy(ProjectID int, privacy string) error {
	query := `UPDATE projects SET privacy = ? WHERE id = ?`
	_, err := models.UserDB.Exec(query, privacy, ProjectID)
	if err != nil {
		return fmt.Errorf("ошибка обновления приватности: %w", err)
	}
	return nil

}


//обновить названия и описания проекта
// обновляем статус приватности проекта
func UpdateProjectInfo(projectID int, name, description string) error {
	query := `UPDATE projects SET name = ?, description = ? WHERE id = ?`
	_, err := models.UserDB.Exec(query, name, description, projectID)
	if err != nil {
		return fmt.Errorf("ошибка обновления приватности: %w", err)
	}
	return nil

}

