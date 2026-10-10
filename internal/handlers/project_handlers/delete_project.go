package project

import (
	projectRepo "HOTA/internal/repositories/project"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type DRODeleteMess struct {
	Message string `json:"message"`
}

// удаление проектов вместе со слотами и заявками
func DeleteProject(w http.ResponseWriter, r *http.Request) {

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

	projectData, err := projectRepo.GetProjectByID(projectID)
	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}

	if userID != projectData.OwnerID {
		http.Error(w, "Удалить проект может только создатель проекта", http.StatusForbidden)
		return
	}

	if projectData.Status == "finished" {
		http.Error(w, "Проект завершён, редакиорвание невозможно", http.StatusConflict)
		return
	}

	//если проект в рабте или завершен, то проект нельзя удалить
	if projectData.Status != "draft" {
		http.Error(w, "Удалить проект можно только в статусе набора разработчиков", http.StatusConflict)
		return
	}

	err = projectRepo.DeleteProjectID(projectData.ID)
	if err != nil {
		fmt.Printf("Ошибка удаления проекта %v", err)
		http.Error(w, "Ошибка удаления проекта", http.StatusInternalServerError)
		return

	}

	res := DRODeleteMess{
		Message: "Проект удален",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}
