package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

//удаляет все стеки пользователя
func DeleteUserStacks(userID int) error {
	query := `DELETE FROM user_stacks WHERE user_id = ?`
	_, err := models.UserDB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("ошибка удаления стеков пользователя: %w", err)
	}
	return nil
}

//добавляет стек пользователю
func AddUserStack(userID, stackID int) error {
	query := `INSERT INTO user_stacks (user_id, stack_id) VALUES (?, ?)`
	_, err := models.UserDB.Exec(query, userID, stackID)
	if err != nil {
		return fmt.Errorf("ошибка добавления стека пользователю: %w", err)
	}
	return nil
}

//проверяет, что все stack_id существуют
func ValidateStacksExist(stackIDs []int) error {
	for _, id := range stackIDs {
		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM stacks WHERE id = ?)`
		err := models.UserDB.QueryRow(query, id).Scan(&exists)
		if err != nil {
			return fmt.Errorf("ошибка проверки стека: %w", err)
		}
		if !exists {
			return fmt.Errorf("стек с id %d не найден", id)
		}
	}
	return nil
}