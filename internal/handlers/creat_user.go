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
	fmt.Println(user)

	//валидиреум пользователя
	errs := service.ValidateStruct(user)
	fmt.Println(errs)
	if errs != nil {
		http.Error(w, "Ошибка валидации", 400)
		status := models.RegistrationResponseError{
			StatusGlobal: "Регистрация не успешна",
			Error:        errs,
		}
		data, _ := json.MarshalIndent(&status,
			" ",
			" ")

		w.Write(data)
		return
	}

	//хешируем + солим пароль
	newpassword, err := service.Hash_password(user.Password)
	if err != nil {
		fmt.Printf("Ошибка хеширования пароля: %v", err)
		statys := models.RegistrationResponseError{
			StatusGlobal: "Ошибка хеширования пароля",
			Error:        err,
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statys)
		return

	}
	user.Password = newpassword

	fmt.Println(user)
	user, err = repositories.AppendUser(user) //сохраняем пользовтеля в БД
	if err != nil {
		status := models.RegistrationResponseError{
			StatusGlobal: "Регистрация не успешна",
			Error:        err,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
		return
	}

	jwtDTO, err := service.CreatJWT(user.ID)
	if err != nil {
		w.WriteHeader(401)
		fmt.Printf("Ошибка создания токена %v", err)
		http.Error(w, "Ошибка создания токена", 400)
		return
	}

	fmt.Println(jwtDTO)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwtDTO)
}
