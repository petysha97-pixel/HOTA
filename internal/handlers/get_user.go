package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// получаем нашего пользователя для отображения данных в главном меню
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
		fmt.Printf("Ошибка получения стеков пользователя %v", err)
		http.Error(w, "Ошибка получения стеков пользователя", http.StatusInternalServerError)
		return
	}
	// Собираем респонс
	reposonse := models.UserResponse{
		ID:       user.ID,
		Nickname: user.Nickname,
		Rolle:    user.Rolle,
		Stack:    stack,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reposonse)
}
