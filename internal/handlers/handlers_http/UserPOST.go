package handlers_http

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(405)
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

		user = repositories.AppendUser(user)

		fmt.Println(user)
	}
}
