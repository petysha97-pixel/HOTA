package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

func DescriptionUser(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}
	fmt.Println(userID)

	var about models.About

	//парсим
	if err := json.NewDecoder(r.Body).Decode(&about); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(about.About) > 200 {
		http.Error(w, "Описание слишком длинное", http.StatusBadRequest)
		return
	}

	//записываем описание о сбее в таблицу
	err = repositories.AboutUser(userID, about.About)
	if err != nil {
		fmt.Printf("ошибка сохранения описания: %v\n", err)
		http.Error(w, "Ошибка сохранения описания", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Описание успешно сохранено",
		"about":   about.About,
	})

}
