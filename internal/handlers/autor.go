package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	var authDTO models.AuotIn

	//парсим
	if err := json.NewDecoder(r.Body).Decode(&authDTO); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//Получаем текущий email и хеш пароля из БД
	var Email string
	var Hashpass string
	var id int

	query := "SELECT id, Email, Password FROM users WHERE Email = ?"
	if err := models.UserDB.QueryRow(query, authDTO.Email).Scan(&id, &Email, &Hashpass); err != nil {
		http.Error(w, "Неверная почта или пароль", http.StatusUnauthorized)
		return
	}

	//сравниванием пароли
	if err := service.CheckPassword(Hashpass, authDTO.Password); err != nil {
		http.Error(w, "Неверная почта или пароль", http.StatusUnauthorized)
		return
	}

	jwtDTO, err := service.CreatJWT(id)
	if err != nil {
		fmt.Printf("Ошибка создания токена %v\n", err)
		http.Error(w, "Ошибка создания токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwtDTO)

}
