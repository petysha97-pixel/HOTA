package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"net/http"
)

//показывает всех разработчиков 
func GetAllUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	users, err := repositories.GETUser()
	if err != nil {
		fmt.Printf("Ошибка всех пользователей %v", err)
		http.Error(w, "Ошибка всех пользователей", http.StatusInternalServerError)
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
