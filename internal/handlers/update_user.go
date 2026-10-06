package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// функция обновления пользователя
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	// Достаем userID из контекста запроса
	// r.Context().Value возвращает тип interface{}, поэтому приводим его к string через .(string)
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	//
	var user models.User

	//парсим
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//валидируем пользователя
	err = service.ValidateUpdataStruct(user, userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Данные не прошли валидацию для обновления", 400)
		return
	}

	updateuser, err := repositories.UpdateUser(user, userID)
	if err != nil {
		fmt.Printf("Ошибка обновления пользователя: %v", err)
		http.Error(w, "Ошибка обновления пользователя", 418)
		return
	}
	if updateuser == nil {
		http.Error(w, "Пользователь не найден", http.StatusBadRequest)
		return
	}

	usersDTO := models.UserResponse{
		ID:       userID,
		Email:    updateuser.Email,
		Nickname: updateuser.Nickname,
		Name:     updateuser.Name,
		Rolle:    updateuser.Rolle,
		Grade:    updateuser.Grade,
		Stack:    []models.Stack{},
	}

	stacks, err := repositories.UpdateUserStackID(userID, user.Stack)
	if err != nil {
		fmt.Printf("Ошибка обновления пользователя: %v", err)
		http.Error(w, "ошибка обновления стеков", http.StatusNotFound)
		return
	}

	for _, stack := range stacks {
		usersDTO.Stack = append(usersDTO.Stack, stack)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usersDTO)

}
