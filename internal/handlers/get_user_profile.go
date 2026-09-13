package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)


// Переходим в профиль другого разраба
func GetUserOtherProfile(w http.ResponseWriter, r *http.Request) {

	userID := r.PathValue("id")
	fmt.Println(userID)

	id, err := strconv.Atoi(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Некорректный id пользователя для перехода на его профиль", http.StatusBadRequest)
		fmt.Printf("Не получилось конвертировать строку в число: %s", userID)
		return
	}
	fmt.Println(id)
	//Провеока пользователя в БД по айди и возвращаем его
	user, err := repositories.GetUsersByID(id)
	if err != nil {
		fmt.Printf("ошибка поиск пользователя в БД: %v", err)
		http.Error(w, "Ошибка поиска пользователя в БД для указания описния о себе", http.StatusBadRequest)
		return
	}
	if user == nil {
		fmt.Printf("ошибка поиск пользователя: %v", err)
		http.Error(w, "Пользователь не найден", http.StatusBadRequest)
		return
	}

	//берём стеки
	stack, err := repositories.GetStacksByUserID(user.ID)
	if err != nil {
		http.Error(w, "Ошибка сканирования стека при передачи всех пользователей", http.StatusInternalServerError)
		return
	}

	responce := models.UserResponse{
		ID:       user.ID,
		Nickname: user.Nickname,
		Rolle:    user.Rolle,
		Stack:    stack,
		About:    user.About,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responce)
}
