package project

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// драфт - воркинг - финиш

// тело запроса: только новый статус
type projectStatusIn struct {
	Status string `json:"status"`
}

// обновляем статус комнаты
func UpdateStatusProject(w http.ResponseWriter, r *http.Request) {
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

	var body projectStatusIn
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if body.Status != "draft" && body.Status != "working" && body.Status != "finished" {
		http.Error(w, "Неизвестный статус проекта", http.StatusBadRequest)
		return
	}

	projectData, err := projectRepo.GetProjectByID(projectID)
	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}

	if userID != projectData.OwnerID {
		http.Error(w, "Смена статуса запрещена", http.StatusForbidden)
		return
	}

	if body.Status == projectData.Status {
		http.Error(w, "Данный статус уже установлен", http.StatusConflict)
		return
	}

	totalSlot, openSlot, _, doneSlot, err := projectRepo.CountSlotStatus(projectID)
	if err != nil {
		fmt.Printf("Ошибка подсчёта слотов: %v\n", err)
		http.Error(w, "Ошибка подсчёта слотов", http.StatusInternalServerError)
		return
	}

	switch projectData.Status {
	case "draft":
		if body.Status != "working" {
			http.Error(w, "Из draft можно только в working", http.StatusConflict)
			return
		}
		if openSlot > 0 {
			http.Error(w, "Перевод запрещён: есть открытые слоты", http.StatusConflict)
			return
		}

	case "working":
		if body.Status == "draft" {
			// создатель снова открыл набор
		} else if body.Status == "finished" {
			if doneSlot != totalSlot {
				http.Error(w, "Перевод запрещён: не все слоты завершены", http.StatusConflict)
				return
			}
		} else {
			http.Error(w, "Из working можно в draft или finished", http.StatusConflict)
			return
		}

	case "finished":
		http.Error(w, "Проект завершён, смена статуса запрещена", http.StatusConflict)
		return

	default:
		http.Error(w, "Неизвестный текущий статус проекта", http.StatusConflict)
		return
	}

	err = projectRepo.UpdateProjectStatus(projectID, body.Status)
	if err != nil {
		http.Error(w, "Ошибка обновления статуса проекта", http.StatusInternalServerError)
		return
	}

	res := models.UpdateStatusRoom{
		Message: "Статус обновлён",
		Status:  body.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

//anykey
//nota-app.dev
