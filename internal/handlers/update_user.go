package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// правка профиля из окна «Профиль»: ник, ФИО, роль, грейд, о себе
// почта меняется ручкой /user/email, стек — ручками /user/stack
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	// Достаем userID из контекста запроса
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	//парсим
	var profile models.UpdateProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//валидируем
	err = service.ValidateUpdataStruct(profile, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = repositories.UpdateUser(profile, userID)
	if err != nil {
		fmt.Printf("Ошибка обновления пользователя: %v\n", err)
		http.Error(w, "Ошибка обновления пользователя", http.StatusInternalServerError)
		return
	}

	// в ответ отдаём обновлённый профиль (без истории участия — она не менялась)
	user := repositories.Get_userdb(userID)
	stacks, err := repositories.GetStacksByUserID(userID)
	if err != nil {
		fmt.Printf("Ошибка получения стеков пользователя: %v\n", err)
		http.Error(w, "Ошибка получения стеков пользователя", http.StatusInternalServerError)
		return
	}

	usersDTO := models.UserResponse{
		ID:       user.ID,
		Nickname: user.Nickname,
		Name:     user.Name,
		Rolle:    user.Rolle,
		Grade:    user.Grade,
		About:    user.About,
		Stack:    stacks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usersDTO)
}
