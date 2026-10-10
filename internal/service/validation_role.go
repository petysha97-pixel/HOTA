package service

import (
	"HOTA/internal/models"
	"errors"
	"fmt"
)

// Роль должна быть в каталоге ролей (таблица roles)
func RoleExists(value any) error {
	rolle, ok := value.(string)
	if !ok {
		return errors.New("роль должна быть строкой")
	}
	if rolle == "" {
		return nil // пустое значение ловит validation.Required
	}

	var exists bool
	err := models.UserDB.QueryRow(`SELECT EXISTS(SELECT 1 FROM roles WHERE name = ?)`, rolle).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки роли: %w", err)
	}
	if !exists {
		return errors.New("такой роли нет в каталоге")
	}

	return nil
}
