package handlers_http

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// регистрация пользователей через POST запросы
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
		fmt.Println(user)

		errs := validation(user)
		if len(errs) > 0 {
			status := models.Registration_fail{
				StatusGlobal: "Регистрация не успешна",
				Error:        errs,
			}
			data, _ := json.MarshalIndent(&status,
				" ",
				" ")

			w.Write(data)
			return
		}
		if len(errs) == 0 {
			status := models.Registration_fail{
				StatusGlobal: "Регистрация успешна",
			}
			data, _ := json.MarshalIndent(&status,
				" ",
				" ")

			w.Write(data)

		}

		user = repositories.AppendUser(user)
		fmt.Println(user)
	}
}

// валидация полей
func validation(user models.User) []string {

	var result []string

	if user.Nickname == "" {
		err := "Пустой никнейм"
		result = append(result, err)
	}
	if user.Role == "" {
		err := "Пустая роль"
		result = append(result, err)
	}
	if len(user.Stack) == 0 {
		err := "Пустой стек"
		result = append(result, err)
	}
	if user.GitHub == "" {
		err := "Пустой гитхаб"
		result = append(result, err)
	}
	if user.Telegram == "" {
		err := "Пустой ТГ"
		result = append(result, err)
	}
	if user.Status == "" {
		err := "Пустой статус"
		result = append(result, err)
	}

	if len(result) > 0 {
		return result
	}

	return nil
}
