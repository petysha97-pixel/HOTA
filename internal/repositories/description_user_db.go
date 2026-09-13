package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

// СОХРАНЯЕТ ИНФОРМАЦИЮ "О СЕБЕ"
func AboutUser(id int, description string) error {
	query := `UPDATE users SET about = ? WHERE id = ?`

	_, err := models.UserDB.Exec(query, description, id)
	if err != nil {
		return fmt.Errorf("ошибка сохранения описания: %w", err)
	}

	return nil
}
