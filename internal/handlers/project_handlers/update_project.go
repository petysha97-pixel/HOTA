package project

import (
	projectRepo "HOTA/internal/repositories/project"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type DTOUpdateProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// редактирования названия и описания проекта
func UpdateProject(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	projectID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Неверный ID проекта", http.StatusBadRequest)
		return
	}

	var body DTOUpdateProject
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if body.Name == "" && body.Description == "" {
		http.Error(w, "Поля пустые, нечего обновлять", http.StatusBadRequest)
		return
	}

	projectData, err := projectRepo.GetProjectByID(projectID)
	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}

	if userID != projectData.OwnerID {
		http.Error(w, "Только создатель может изменить проект", http.StatusForbidden)
		return
	}

	if projectData.Status == "finished" {
		http.Error(w, "Проект завершён, редакиорвание невозможно", http.StatusConflict)
		return
	}

	// Обновление названия и описания у проекиа
	err = projectRepo.UpdateProjectInfo(projectData.ID, body.Name, body.Description)
	if err != nil {
		fmt.Printf("ошибка обновления проекта %v", err)
		http.Error(w, "Ошибка обновления проекта", http.StatusInternalServerError)
		return
	}

	newproject, err := projectRepo.GetProjectByID(projectID)
	if err != nil {
		http.Error(w, "Ошибка получения проекта", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newproject)

}
