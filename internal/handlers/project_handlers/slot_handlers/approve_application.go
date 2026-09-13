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

// принять заявку на слот
func ApproveApplication(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Только владелец проекта может утверждать заявки", http.StatusForbidden)
		return
	}

	// проверяем, что слот ещё открыт
	if slotData.Status != "open" {
		http.Error(w, "Слот уже закрыт", http.StatusBadRequest)
		return
	}

	// проверяем, что заявка существует и на рассмотрении
	app, err := slotRepo.GetUserCheckingApplications(slotID, applicantID)
	if err != nil {
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}
	if app.Status != "pending" {
		http.Error(w, "Заявка уже обработана", http.StatusBadRequest)
		return
	}

	// обновляем статус заявки
	err = slotRepo.UpdateApplicationStatus(app.ID, "approved")
	if err != nil {
		http.Error(w, "Ошибка обновления заявки", http.StatusInternalServerError)
		return
	}

	// закрываем слот
	_, err = slotRepo.UpdateSlotByID(slotID, "closed")
	if err != nil {
		http.Error(w, "Ошибка закрытия слота", http.StatusInternalServerError)
		return
	}

	// отклоняем все остальные заявки на этот слот
	err = slotRepo.RejectOtherApplications(slotID, app.ID)
	if err != nil {
		http.Error(w, "Ошибка отклонения остальных заявок", http.StatusInternalServerError)
		return
	}

	// ответ
	res := models.ApplicationResponse{
		Message: "Заявка принята, слот закрыт",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
