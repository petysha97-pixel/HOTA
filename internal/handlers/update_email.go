package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

func UpdateEmail(w http.ResponseWriter, r *http.Request) {

	// Достаем userID из контекста запроса
	// r.Context().Value возвращает тип interface{}, поэтому приводим его к string через .(string)
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	var emailData models.UpdateEmail

	//парсим
	if err := json.NewDecoder(r.Body).Decode(&emailData); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = service.ValidateUpdateEmail(userID, emailData.Email, emailData.Password)
	if err != nil {
		fmt.Printf("Ошибка обновления email: %v", err)
		http.Error(w, "ошибка обновления email", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Email успешно обновлен")
}
