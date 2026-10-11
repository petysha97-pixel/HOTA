package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// вкладка «Разработчики»: список и поиск
// GET /users?query=go&rolle=Backend&grade=Middle&stackID=1,5
// все параметры необязательные, без них — все разработчики (кроме меня)
func GetUsers(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// собираем фильтры из адресной строки
	filter := models.UserFilter{
		Query: strings.TrimSpace(r.URL.Query().Get("query")),
		Rolle: r.URL.Query().Get("rolle"),
		Grade: r.URL.Query().Get("grade"),
	}

	// stackID приходит строкой "1,5,7" — переводим в список чисел
	stackText := r.URL.Query().Get("stackID")
	if stackText != "" {
		for _, part := range strings.Split(stackText, ",") {
			id, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				http.Error(w, "stackID — это id технологий через запятую", http.StatusBadRequest)
				return
			}
			filter.StackID = append(filter.StackID, id)
		}
	}

	if err := service.ValidateUserFilter(filter); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, err := repositories.GetUsers(filter, userID)
	if err != nil {
		fmt.Printf("Ошибка получения разработчиков %v\n", err)
		http.Error(w, "Ошибка получения разработчиков", http.StatusInternalServerError)
		return
	}

	usersDTO := []models.UserResponse{}
	for _, user := range users {
		// берём стек вместе с описанием опыта (нужен для карточки разработчика)
		stack, err := repositories.GetStacksByUserID(user.ID)
		if err != nil {
			fmt.Printf("Ошибка получения стека разработчика %v\n", err)
			http.Error(w, "Ошибка получения стека разработчика", http.StatusInternalServerError)
			return
		}

		// текстовый поиск: ник, имя, роль или название технологии
		if !textFound(filter.Query, user, stack) {
			continue
		}

		usersDTO = append(usersDTO, models.UserResponse{
			ID:       user.ID,
			Nickname: user.Nickname,
			Name:     user.Name,
			Rolle:    user.Rolle,
			Grade:    user.Grade,
			About:    user.About,
			Stack:    stack,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersDTO)
}

// есть ли текст поиска в нике, имени, роли или технологиях разработчика
// ищем без учёта больших и маленьких букв (работает и для русских)
func textFound(text string, user models.User, stack []models.Stack) bool {
	if text == "" {
		return true // текста нет — подходит любой
	}

	vse := []string{user.Nickname, user.Name, user.Rolle}
	for _, s := range stack {
		vse = append(vse, s.Name)
	}

	stroka := strings.ToLower(strings.Join(vse, " "))
	return strings.Contains(stroka, strings.ToLower(text))
}
