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

// удалить одну технологию из стека: последнюю удалить нельзя
func DeleteStack(w http.ResponseWriter, r *http.Request) {
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

	err = repositories.DeleteStackUser(userID, stack.StackID)
	if err != nil {
		if errors.Is(err, repositories.ErrStackNotFound) {
			http.Error(w, "Стек не найден", http.StatusNotFound)
			return
		}
		// последнюю технологию удалить нельзя
		if errors.Is(err, repositories.ErrStackLast) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		fmt.Printf("Ошибка удаления стека %v\n", err)
		http.Error(w, "Ошибка удаления стека", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Стек успешно удален")

}
