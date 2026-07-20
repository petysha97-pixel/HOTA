package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// регистрация +валидация пользователей через POST запросы
func NewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(405)
		return
	}

	if r.Method == "POST" {
		var user models.User

		body, err := io.ReadAll(r.Body)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Не верные данные 1 ")
			return
		}

		err = json.Unmarshal(body, &user)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Не верные данные 2")
			return
		}

		//валидиреум пользователя

		errs := repositories.ValidateStruct(user)
		fmt.Println(errs)
		if errs != nil {
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

		//После валидации хешируем пооль + солим
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "ошибка хеширования пароля")
			return
		}
		user.Password = string(hash)

		user, err = repositories.AppendUser(user) //сохраняем пользовтеля в БД
		fmt.Println(err)
		if err != nil {
			status := models.RegistrationResponseError{
				StatusGlobal: "Регистрация не успешна",
				Error:        err,
			}

			data, _ := json.MarshalIndent(&status,
				" ",
				" ")

			w.Write(data)
			return
		} else {
			status := models.RegistrationResponseError{
				StatusGlobal: "Регистрация успешна",
				ID:           user.ID,
			}
			data, _ := json.MarshalIndent(&status,
				" ",
				" ")

			w.Write(data)
			fmt.Println(string(data))
		}

		fmt.Println(user)
	}
}
