package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

func AddUserStack(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	var stack models.StackUser

	//парсим
	if err := json.NewDecoder(r.Body).Decode(&stack); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = repositories.AddStackUser(userID, stack.IDStack)
	if err != nil {
		fmt.Printf("Ошибка добавления стека %v", err)
		http.Error(w, "Ошибка добавления стека", 505)

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Стек успешно добавлен")

}
