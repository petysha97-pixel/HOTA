package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	err = repositories.DeleteUser(userID)
	if err != nil {
		fmt.Printf("ошибка удаления пользователя в БД: %v", err)
		http.Error(w, "Ошибка удаления пользователя", http.StatusBadRequest)
	}

	otvet := models.RegistrationResponseError{
		ID:           userID,
		StatusGlobal: "Пользователь удален успешно",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(otvet)

}
