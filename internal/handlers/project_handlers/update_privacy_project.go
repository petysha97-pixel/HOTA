package project

import (
	projectRepo "HOTA/internal/repositories/project"
	"HOTA/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type dtoPrivateProject struct {
	Privacy string `json:"privacy"`
}

type UpdatePrivacyProject struct {
	Message string `json:"message"`
	Privacy string `json:"privacy"`
}

// смена приватности на public или private
func UpdatePrivateProject(w http.ResponseWriter, r *http.Request) {

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

	var body dtoPrivateProject
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if body.Privacy != "public" && body.Privacy != "private" {
		http.Error(w, "Не верный статцс приватности", http.StatusBadRequest)
		return
	}

	projectData, err := projectRepo.GetProjectByID(projectID)
	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}

	if userID != projectData.OwnerID {
		http.Error(w, "Только создатель может изменить приватность проекта", http.StatusForbidden)
		return
	}

	if body.Privacy == projectData.Privacy {
		http.Error(w, "Эта приватность уже установлена", http.StatusConflict)
		return
	}

	if projectData.Status == "finished" {
		http.Error(w, "Проект завершён, приватность менять нельзя", http.StatusConflict)
		return
	}

	err = projectRepo.UpdateProjectPrivacy(projectData.ID, body.Privacy)
	if err != nil {
		http.Error(w, "ошибка обновления приватности", http.StatusInternalServerError)
		return
	}

	res := UpdatePrivacyProject{
		Message: "Приватность обновлена",
		Privacy: body.Privacy,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
