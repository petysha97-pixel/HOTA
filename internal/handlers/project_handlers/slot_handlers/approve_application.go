package slot

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	slotRepo "HOTA/internal/repositories/project/slot"
	"HOTA/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

// принять заявку на слот.
// заявка только назначает разработчика: статус слота создатель меняет отдельно
func ApproveApplication(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	slotID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Неверный ID слота", http.StatusBadRequest)
		return
	}
	applicantID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	slotData, err := slotRepo.GetSlotByID(slotID)
	if err != nil {
		http.Error(w, "Слот не найден", http.StatusNotFound)
		return
	}

	projectData, err := projectRepo.GetProjectByID(slotData.ProjectID)
	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}
	if projectData.OwnerID != userID {
		http.Error(w, "Только владелец проекта может утверждать заявки", http.StatusForbidden)
		return
	}

	if projectData.Status == "finished" {
		http.Error(w, "Проект завершён, заявки менять нельзя", http.StatusConflict)
		return
	}

	if slotData.Status != "open" {
		http.Error(w, "Слот уже закрыт", http.StatusConflict)
		return
	}

	// на слоте не должно быть уже утверждённого разработчика
	hasApproved, err := slotRepo.ApprovedAplecation(slotID)
	if err != nil {
		http.Error(w, "Ошибка проверки заявки", http.StatusInternalServerError)
		return
	}
	if hasApproved {
		http.Error(w, "На слот уже утверждён разработчик", http.StatusConflict)
		return
	}

	app, err := slotRepo.GetUserCheckingApplications(slotID, applicantID)
	if err != nil {
		http.Error(w, "Ошибка поиска заявки", http.StatusInternalServerError)
		return
	}
	if app == nil {
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}

	err = slotRepo.UpdateApplicationStatus(app.ID, "approved")
	if err != nil {
		http.Error(w, "Ошибка обновления заявки", http.StatusInternalServerError)
		return
	}

	err = slotRepo.RejectOtherApplications(slotID, app.ID)
	if err != nil {
		http.Error(w, "Ошибка отклонения остальных заявок", http.StatusInternalServerError)
		return
	}

	// слот остаётся open: в close его переводит создатель отдельным запросом
	res := models.ApplicationResponse{
		Message: "Заявка принята, разработчик назначен на слот",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
