package handlers

import (
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"net/http"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet{
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("languages")
	if query == "" {
		fmt.Fprintf(w, "Пустой запрос")
		return
	}

	users, err := repositories.SearchUser(query)

	if err != nil {
		http.Error(w, "Ошибка при поиске пользователей", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
