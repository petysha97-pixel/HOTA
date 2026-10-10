package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

// DeleteStackUser удаляет одну технологию из стека пользователя
// последнюю удалить нельзя: стек обязателен
func DeleteStackUser(userID int, stackID int) error {

	// есть ли такая технология у пользователя
	var ex bool
	err := models.UserDB.QueryRow(`SELECT EXISTS(SELECT 1 FROM user_stacks WHERE user_id = ? AND stack_id = ?)`, userID, stackID).Scan(&ex)
	if err != nil {
		return fmt.Errorf("ошибка проверки стека пользователя: %w", err)
	}
	if !ex {
		return ErrStackNotFound
	}

	// последнюю не удаляем
	count, err := CountUserStacks(userID)
	if err != nil {
		return err
	}
	if count <= 1 {
		return ErrStackLast
	}

	_, err = models.UserDB.Exec(`DELETE FROM user_stacks WHERE user_id = ? AND stack_id = ?`, userID, stackID)
	if err != nil {
		return fmt.Errorf("ошибка удаления стека пользователя: %w", err)
	}

	return nil
}
