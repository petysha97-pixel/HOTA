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
		http.Error(w, "Данный email не зарегистрирован", http.StatusBadRequest)
		return
	}

	//сравниванием пароли
	if err := service.CheckPassword(Hashpass, authDTO.Password); err != nil {
		w.WriteHeader(401)
		fmt.Fprintf(w, "Пароли не совпадают")
		return
	}

	jwtDTO, err := service.CreatJWT(id)
	if err != nil {
		w.WriteHeader(401)
		fmt.Printf("Ошибка создания токена %v", err)
		http.Error(w, "Ошибка создания токена", 400)
		return
	}

	fmt.Println(jwtDTO)
	w.Header().Set("Context-Type", "application/json")
	json.NewEncoder(w).Encode(jwtDTO)

}
