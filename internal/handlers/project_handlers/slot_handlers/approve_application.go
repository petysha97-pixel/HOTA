package slot

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	slotRepo "HOTA/internal/repositories/project/slot"
	"HOTA/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// Утверждаем кандидата на знакомстве и закрываем слот.
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

	if slotData.Status != "review" {
		http.Error(w, "Сначала возьмите кандидата на знакомство", http.StatusConflict)
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

	// Назначаем исполнителя, закрываем слот и отклоняем оставшуюся очередь.
	err = slotRepo.ApproveSlotApplication(slotID, applicantID, userID)
	if errors.Is(err, slotRepo.ErrApprovalConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Ошибка утверждения заявки", http.StatusInternalServerError)
		return
	}

	// Отвечаем после закрытия слота
	res := models.ApplicationResponse{
		Message: "Исполнитель назначен. Набор закрыт.",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
