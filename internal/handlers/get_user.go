package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

//получаем нашего пользователя для отображения данных
func GetUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный id", http.StatusNotFound)
		fmt.Printf("Не получилось конвертировать строку в число: %s", idStr)
		return
	}

	user := repositories.Get_userdb(id)

	respons := models.UserResponse{
		ID:       user.ID,
		Nickname: user.Nickname,
		Rolle:    user.Rolle,
		Stack:    user.Stack,
	}

	json.NewEncoder(w).Encode(respons)
}
