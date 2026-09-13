package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

func UpdatePassword(w http.ResponseWriter, r *http.Request) {

	// Достаем userID из контекста запроса
	// r.Context().Value возвращает тип interface{}, поэтому приводим его к string через .(string)
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	//
	var pass models.UpdatePasswors

	//парсим
	if err := json.NewDecoder(r.Body).Decode(&pass); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = service.ValidateUpdatePass(userID, pass.OldPassword, pass.NewPassword)
	if err != nil {
		fmt.Printf("Ошибка обновления пароля: %v", err)
		http.Error(w, "ошибка обновления пароля", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Пароль успешно обновлен")
}
