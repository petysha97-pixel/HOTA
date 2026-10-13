package project

import (
	projectRepo "HOTA/internal/repositories/project"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// левая панель «Проекты»: «Создал» и «Участвую»
// GET /my/projects
func GetMyProjects(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	response, err := projectRepo.GetMyProjects(userID)
	if err != nil {
		fmt.Printf("Ошибка получения моих проектов %v\n", err)
		http.Error(w, "Ошибка получения моих проектов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
