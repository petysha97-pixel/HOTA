package service

import (
	"HOTA/internal/models"
	"errors"
	"fmt"
)

// проверяем пароль пользователя перед необратимым действием (например, удалением аккаунта)
func ConfirmPassword(id int, password string) error {
	if password == "" {
		return errors.New("пароль не может быть пустым")
	}

	var storedHash string
	err := models.UserDB.QueryRow("SELECT Password FROM users WHERE id = ?", id).Scan(&storedHash)
	if err != nil {
		return fmt.Errorf("ошибка получения данных пользователя: %w", err)
	}

	if err := CheckPassword(storedHash, password); err != nil {
		return errors.New("неверный пароль")
	}

	return nil
}
