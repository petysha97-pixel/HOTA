package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

// DeleteStackUser удаляет стек пользователя по userID и stackID
func DeleteStackUser(userID int, stackID int) error {
	query := `DELETE FROM user_stacks WHERE user_id = ? AND stack_id = ?`

	result, err := models.UserDB.Exec(query, userID, stackID)
	if err != nil {
		return fmt.Errorf("delete stack user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("стек не найден или не принадлежит пользователю")
	}

	return nil
}
