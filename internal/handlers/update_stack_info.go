package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// PATCH /user/stack/{id} — изменить описание опыта у одной технологии
// {id} — id технологии из каталога stacks
// тело: {"description": "Три года на бэкенде..."}
func UpdateStackInfo(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	stackID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Неверный ID стека", http.StatusBadRequest)
		return
	}

	//парсим
	var info models.StackChange
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// описание до 300 символов
	if err := service.ValidateStackInfo(info); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = repositories.UpdateStackInfo(userID, stackID, info.Description)
	if err != nil {
		if errors.Is(err, repositories.ErrStackNotFound) {
			http.Error(w, "Этой технологии нет в вашем стеке", http.StatusNotFound)
			return
		}
		fmt.Printf("Ошибка изменения стека: %v\n", err)
		http.Error(w, "Ошибка изменения стека", http.StatusInternalServerError)
		return
	}

	// в ответ — весь стек пользователя, как его показывает профиль
	stacks, err := repositories.GetStacksByUserID(userID)
	if err != nil {
		fmt.Printf("Ошибка получения стеков пользователя: %v\n", err)
		http.Error(w, "Ошибка получения стеков пользователя", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stacks)
}
