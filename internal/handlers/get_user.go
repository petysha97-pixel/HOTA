package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// свой профиль: вся страница кроме почты и пароля
// ник, имя, роль, грейд, о себе, стек (с уровнем и описанием), история участия
func GetUser(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Получаем юзера
	user := repositories.Get_userdb(userID)

	// Получаем его стеки
	stack, err := repositories.GetStacksByUserID(userID)
	if err != nil {
		fmt.Printf("Ошибка получения стеков пользователя %v\n", err)
		http.Error(w, "Ошибка получения стеков пользователя", http.StatusInternalServerError)
		return
	}

	// история участия: в своём профиле видны и приватные проекты
	history, err := repositories.GetUserHistory(userID, false)
	if err != nil {
		fmt.Printf("Ошибка получения истории участия %v\n", err)
		http.Error(w, "Ошибка получения истории участия", http.StatusInternalServerError)
		return
	}

	// Собираем респонс
	reposonse := models.UserResponse{
		ID:       user.ID,
		Nickname: user.Nickname,
		Name:     user.Name,
		Rolle:    user.Rolle,
		Grade:    user.Grade,
		About:    user.About,
		Stack:    stack,
		History:  history,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reposonse)
}
