package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"

	// "crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
)

// регистрация + валидация пользователей через POST запросы
func NewUser(w http.ResponseWriter, r *http.Request) {

	var user models.User

	//парсим
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//валидиреум пользователя
	errs := service.ValidateStruct(user)
	fmt.Println(errs)
	if errs != nil {
		status := models.RegistrationResponseError{
			StatusGlobal: "Регистрация не успешна",
			Error:        errs,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(status)
		return
	}

	//хешируем + солим пароль
	newpassword, err := service.Hash_password(user.Password)
	if err != nil {
		fmt.Printf("Ошибка хеширования пароля: %v\n", err)
		http.Error(w, "Ошибка хеширования пароля", http.StatusInternalServerError)
		return

	}
	user.Password = newpassword

	user, err = repositories.AppendUser(user) //сохраняем пользовтеля в БД
	if err != nil {
		fmt.Printf("Ошибка сохранения пользователя: %v\n", err)
		http.Error(w, "Регистрация не успешна", http.StatusInternalServerError)
		return
	}

	jwtDTO, err := service.CreatJWT(user.ID)
	if err != nil {
		fmt.Printf("Ошибка создания токена %v\n", err)
		http.Error(w, "Ошибка создания токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwtDTO)
}
