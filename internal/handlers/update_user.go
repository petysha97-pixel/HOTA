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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateuser, err := repositories.UpdateUser(user, userID)
	if err != nil {
		fmt.Printf("Ошибка обновления пользователя: %v", err)
		http.Error(w, "Ошибка обновления пользователя", http.StatusInternalServerError)
		return
	}
	if updateuser == nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
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

	// стек уже обновлён в той же транзакции, берём его для ответа
	stacks, err := repositories.GetStacksByUserID(userID)
	if err != nil {
		fmt.Printf("Ошибка получения стеков пользователя: %v\n", err)
		http.Error(w, "Ошибка получения стеков пользователя", http.StatusInternalServerError)
		return
	}

	for _, stack := range stacks {
		usersDTO.Stack = append(usersDTO.Stack, stack)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usersDTO)

}
