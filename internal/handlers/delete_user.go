package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// тело запроса на удаление аккаунта: только пароль для подтверждения
type deleteUserIn struct {
	Password string `json:"password"`
}

// удаление аккаунта: нужен пароль, удаляются проекты пользователя, его заявки и стек
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	//парсим
	var body deleteUserIn
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// без верного пароля аккаунт не удаляем
	if err := service.ConfirmPassword(userID, body.Password); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = repositories.DeleteUser(userID)
	if err != nil {
		fmt.Printf("ошибка удаления пользователя в БД: %v\n", err)
		http.Error(w, "Ошибка удаления пользователя", http.StatusInternalServerError)
		return
	}

	otvet := models.RegistrationResponseError{
		ID:           userID,
		StatusGlobal: "Пользователь удален успешно",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(otvet)

}
