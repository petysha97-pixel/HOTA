package service

import (
	"HOTA/internal/repositories"
	"fmt"
	"net/http"
	"strconv"
)

// берет айди из котеста и преобразует в инт и проверяет, есть юзер с айди он в БД (частое использование)
func ContextUserIDValid(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(userIDKey).(string) 
	if !ok {
		return 0, fmt.Errorf("ID пользователя не найден")
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, fmt.Errorf("некорректный ID: %s", userID)
	}

	usr, err := repositories.GetUsersByID(id)
	if err != nil {
		return 0, fmt.Errorf("ошибка поиска: %w", err)
	}
	if usr == nil {
		return 0, fmt.Errorf("пользователь не найден")
	}

	return id, nil
}
