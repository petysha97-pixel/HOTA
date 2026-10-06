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

	if len(req.StackID) < 1 || len(req.StackID) > 6 {
		http.Error(w, "Стек: от 1 до 6 технологий", http.StatusBadRequest)
		return
	}

	//проверяем, что стеки существуют и не повторяются
	if err := repositories.ValidateStacksExist(req.StackID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//заменяем стек одной транзакцией: при ошибке остаётся старый
	if err := repositories.ReplaceUserStacks(userID, req.StackID); err != nil {
		fmt.Printf("Ошибка обновления стеков: %v\n", err)
		http.Error(w, "Ошибка обновления стеков", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Стеки успешно обновлены",
		"stackID": req.StackID,
	})
}
