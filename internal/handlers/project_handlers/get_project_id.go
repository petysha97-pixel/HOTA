package project

import (
	projectRepo "HOTA/internal/repositories/project"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// страница проекта: проект, ник автора и слоты со стеком
// GET /project/{id}
// приватный проект видят только создатель и исполнители его слотов
func GetProject(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// вытаскиваем ID из пути /project/123
	projectID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Неверный ID проекта", http.StatusBadRequest)
		return
	}

	project, err := projectRepo.GetProjectCard(projectID)
	if err != nil {
		fmt.Printf("Ошибка получения проекта %v\n", err)
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}

	// приватный: пускаем создателя и исполнителей
	if project.Privacy == "private" {
		mozhno := project.OwnerID == userID
		for _, slot := range project.Slots {
			if slot.UserID != nil && *slot.UserID == userID {
				mozhno = true
			}
		}
		if !mozhno {
			http.Error(w, "Доступ запрещён", http.StatusForbidden)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}
