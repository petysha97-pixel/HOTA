package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// добавить одну технологию в стек (на странице профиля): не больше 6
// тело: {"stackID": 1, "description": "..."} — description можно не передавать
func AddUserStack(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	//парсим
	var stack models.StackUser
	if err := json.NewDecoder(r.Body).Decode(&stack); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// описание (если передано) проверяем
	info := models.StackChange{Description: stack.Description}
	if err := service.ValidateStackInfo(info); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = repositories.AddStackUser(userID, stack.StackID, stack.Description)
	if err != nil {
		if errors.Is(err, repositories.ErrStackNotFound) {
			http.Error(w, "Стек не найден", http.StatusNotFound)
			return
		}
		if errors.Is(err, repositories.ErrStackAlreadyAdded) || errors.Is(err, repositories.ErrStackLimit) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		fmt.Printf("Ошибка добавления стека %v\n", err)
		http.Error(w, "Ошибка добавления стека", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Стек успешно добавлен")

}
