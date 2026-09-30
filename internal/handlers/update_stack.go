package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)


func UpdateStack(w http.ResponseWriter, r *http.Request) {
	// 1. Проверяем пользователя
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}
	var req models.UpdateStacks
	//парсим
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//проверяем, что стеки существуют
	if err := repositories.ValidateStacksExist(req.StackID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Удаляем старые стеки
	if err := repositories.DeleteUserStacks(userID); err != nil {
		fmt.Printf("Ошибка удаления стеков: %v\n", err)
		http.Error(w, "Ошибка обновления стеков", http.StatusInternalServerError)
		return
	}

	//Добавляем новые стеки
	for _, stackID := range req.StackID {
		if err := repositories.AddUserStack(userID, stackID); err != nil {
			fmt.Printf("Ошибка добавления стека: %v\n", err)
			http.Error(w, "Ошибка обновления стеков", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Стеки успешно обновлены",
		"stack_ids": req.StackID,
	})
}
