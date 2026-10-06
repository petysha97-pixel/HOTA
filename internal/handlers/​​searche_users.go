package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"net/http"
)

func SearcheUsers(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("query")
	limit := r.URL.Query().Get("limit")

	if limit == "" {
		limit = "15"

	}

	if query == "" {
		http.Error(w, "Параметр query пустой", http.StatusBadRequest) //400 не верный запрос
		return
	}

	users, err := repositories.SearcheUsersBD(query, limit)
	if err != nil {
		fmt.Printf("Ошибка поиска пользователей %v", err)
		http.Error(w, "Ошибка поиска пользователей", http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		http.Error(w, "По запросу ничего не найдено", http.StatusInternalServerError)
		return
	}

	var usersDTO = make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		//берём стеки
		stack, err := repositories.GetStacksByUserID(user.ID)
		if err != nil {
			http.Error(w, "Ошибка сканирования стека при передачи всех пользователей", http.StatusInternalServerError)
			return
		}
		//заполняем нашу ДТО структу
		usersDTO = append(usersDTO, models.UserResponse{
			ID:       user.ID,
			Nickname: user.Nickname,
			Name:     user.Name,
			Rolle:    user.Rolle,
			Grade:    user.Grade,
			Stack:    stack,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersDTO)

}
