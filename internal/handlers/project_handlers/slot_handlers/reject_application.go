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

// отклонить заявку на слот
func RejectApplication(w http.ResponseWriter, r *http.Request) {
	// проверяем пользователя
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// достаём ID слота и пользователя из URL
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

	// проверяем, что слот существует
	slotData, err := slotRepo.GetSlotByID(slotID)
	if err != nil {
		http.Error(w, "Слот не найден", http.StatusNotFound)
		return
	}

	// проверяем, что пользователь — владелец проекта
	projectData, err := projectRepo.GetProjectByID(slotData.ProjectID)
	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}
	if projectData.OwnerID != userID {
		http.Error(w, "Только владелец проекта может отклонять заявки", http.StatusForbidden)
		return
	}

	if projectData.Status == "finished" {
		http.Error(w, "Проект завершён, заявки менять нельзя", http.StatusConflict)
		return
	}

	// проверяем, что заявка существует и на рассмотрении
	app, err := slotRepo.GetUserCheckingApplications(slotID, applicantID)
	if err != nil {
		http.Error(w, "Ошибка поиска заявки", http.StatusInternalServerError)
		return
	}
	if app == nil {
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}

	// Отклоняем заявку и снимаем кандидата со знакомства.
	err = slotRepo.RejectSlotApplication(slotID, applicantID, userID)
	if errors.Is(err, slotRepo.ErrRejectConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Ошибка отклонения заявки", http.StatusInternalServerError)
		return
	}

	// ответ
	res := models.ApplicationResponse{
		Message: "Заявка отклонена",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
