package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Переходим в профиль другого разраба: та же страница, что и свой профиль, но без почты
// в истории участия только публичные проекты
func GetUserOtherProfile(w http.ResponseWriter, r *http.Request) {

	userID := r.PathValue("id")

	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Некорректный id пользователя для перехода на его профиль", http.StatusBadRequest)
		return
	}

	//Провеока пользователя в БД по айди и возвращаем его
	user, err := repositories.GetUsersByID(id)
	if err != nil || user == nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	//берём стеки
	stack, err := repositories.GetStacksByUserID(user.ID)
	if err != nil {
		fmt.Printf("Ошибка получения стеков пользователя %v\n", err)
		http.Error(w, "Ошибка получения стеков пользователя", http.StatusInternalServerError)
		return
	}

	// история участия: чужим показываем только публичные проекты
	history, err := repositories.GetUserHistory(user.ID, true)
	if err != nil {
		fmt.Printf("Ошибка получения истории участия %v\n", err)
		http.Error(w, "Ошибка получения истории участия", http.StatusInternalServerError)
		return
	}

	responce := models.UserResponse{
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
	json.NewEncoder(w).Encode(responce)
}
