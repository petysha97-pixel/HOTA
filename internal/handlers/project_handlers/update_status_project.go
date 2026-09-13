package project

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	"HOTA/internal/service"
	"encoding/json"
	"net/http"
)

// обновляем статус комнаты
func UpdateStatusProject(w http.ResponseWriter, r *http.Request) {
	// берём из контекста + валидация
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// парсим тело
	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// достаём комнату из БД по ID из тела
	roomDa, err := projectRepo.GetProjectByID(project.ID)
	if err != nil {
		http.Error(w, "Комната не найдена", http.StatusNotFound)
		return
	}
	if roomDa.OwnerID != userID {
		http.Error(w, "Только владелец может поменять статус комнаты", http.StatusForbidden)
		return
	}

	newStatus := project.Status

	// все статусы
	if newStatus != "draft" && newStatus != "working" && newStatus != "finished" {
		http.Error(w, "Допустимые статусы: draft, working, finished", http.StatusBadRequest)
		return
	}

	// нельзя из finished обратно
	if roomDa.Status == "finished" && newStatus != "finished" {
		http.Error(w, "Нельзя изменить статус завершённого проекта", http.StatusBadRequest)
		return
	}

	// при переход на ворк все слоты должныть быть закрыты
	if newStatus == "working" {
		slots, err := projectRepo.GetSlotsByProjectID(project.ID)
		if err != nil {
			http.Error(w, "Ошибка проверки слотов", http.StatusInternalServerError)
			return
		}
		for _, slot := range slots {
			if slot.Status != "closed" {
				http.Error(w, "Нельзя перейти в working: есть открытые слоты", http.StatusBadRequest)
				return
			}
		}
	}

	// обновляем статус
	err = projectRepo.UpdateProjectStatus(project.ID, newStatus)
	if err != nil {
		http.Error(w, "Ошибка обновления статуса", http.StatusInternalServerError)
		return
	}

	// ответ
	req := models.UpdateStatusRoom{
		Message: "Статус обновлён",
		Status:  newStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}

//anykey
//nota-app.dev
