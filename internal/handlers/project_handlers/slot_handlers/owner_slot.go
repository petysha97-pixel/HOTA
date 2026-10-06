package slot

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	slotRepo "HOTA/internal/repositories/project/slot"
	"HOTA/internal/service"
	"net/http"
	"strconv"
)

// общий доступ к слоту для создателя проекта: /slot/{id}
// при ошибке сам пишет ответ и возвращает ok = false
func ownerSlot(w http.ResponseWriter, r *http.Request) (*models.Slot, *models.Project, bool) {
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return nil, nil, false
	}

	slotID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || slotID <= 0 {
		http.Error(w, "Неверный ID слота", http.StatusBadRequest)
		return nil, nil, false
	}

	slotData, err := slotRepo.GetSlotByID(slotID)
	if err != nil {
		http.Error(w, "Слот не найден", http.StatusNotFound)
		return nil, nil, false
	}

	projectData, err := projectRepo.GetProjectByID(slotData.ProjectID)
	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return nil, nil, false
	}

	if projectData.OwnerID != userID {
		http.Error(w, "Только создатель проекта может менять слот", http.StatusForbidden)
		return nil, nil, false
	}

	if projectData.Status == "finished" {
		http.Error(w, "Проект завершён, слоты менять нельзя", http.StatusConflict)
		return nil, nil, false
	}

	return slotData, projectData, true
}
