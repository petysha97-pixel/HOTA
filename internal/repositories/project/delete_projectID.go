package project

import (
	"HOTA/internal/models"
	"fmt"
)

//каскадное удаление проекта
func DeleteProjectID(prID int) error {
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() // откатится, если не закоммитим

	// 1. Удаляем заявки на слоты этого проекта
	_, err = tx.Exec(`
		DELETE FROM applicationsSlot
		WHERE slot_id IN (
			SELECT id FROM slot_id WHERE project_id = ?
		)
	`, prID)
	if err != nil {
		return fmt.Errorf("delete applications: %w", err)
	}

	// 2. Удаляем слоты проекта
	_, err = tx.Exec(`DELETE FROM slots WHERE project_id = ?`, prID)
	if err != nil {
		return fmt.Errorf("delete slots: %w", err)
	}

	// 3. Удаляем сам проект
	res, err := tx.Exec(`DELETE FROM projects WHERE id = ?`, prID)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("project %d not found", prID)
	}

	return tx.Commit()
}
