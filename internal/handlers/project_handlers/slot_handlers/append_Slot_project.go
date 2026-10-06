package slot

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	"encoding/json"
	"fmt"
	"strconv"

	"HOTA/internal/service"
	"net/http"
)

// POST /project/{id}/slot — добавить слот в созданный проект
func CreateSlotInProject(w http.ResponseWriter, r *http.Request) {
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

	var body models.DTOSlot
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// роль из списка, стек от 1 до 6
	if err := service.ValidateSlotStruct(body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	projectData, err := projectRepo.GetProjectByID(projectID)
	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}
	if userID != projectData.OwnerID {
		http.Error(w, "Только создатель может добавлять слоты", http.StatusForbidden)
		return
	}
	if projectData.Status == "finished" {
		http.Error(w, "Проект завершён, слоты добавлять нельзя", http.StatusConflict)
		return
	}

	totalSlot, _, _, _, err := projectRepo.CountSlotStatus(projectID)
	if err != nil {
		http.Error(w, "Ошибка подсчёта слотов", http.StatusInternalServerError)
		return
	}
	if totalSlot >= 6 {
		http.Error(w, "В проекте уже 6 слотов", http.StatusConflict)
		return
	}

	slotModel := &models.Slot{
		ProjectID: projectID,
		Rolle:     body.Rolle,
		StackID:   body.Stack,
		Status:    "open", // новый слот всегда открыт
	}
	if err = projectRepo.CreateSlot(slotModel); err != nil {
		fmt.Printf("Ошибка создания слота %v\n", err)
		http.Error(w, "Ошибка создания слота", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(slotModel)
}
