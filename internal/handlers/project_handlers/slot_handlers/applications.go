package slot

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	slotRepo "HOTA/internal/repositories/project/slot"
	"encoding/json"
	"strconv"

	"HOTA/internal/service"
	"net/http"
)

// создание заявки на слот
func ApplicationsSlot(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	slotID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Некорректный id слота", http.StatusBadRequest)
		return
	}

	// проверка слота в БД
	slotData, err := slotRepo.GetSlotByID(slotID)
	if err != nil {
		http.Error(w, "Слот не найден", http.StatusNotFound)
		return
	}

	if slotData.Status != "open" {
		http.Error(w, "Слот закрыт", http.StatusConflict)
		return
	}

	// проект слота: нужен для проверки владельца и статуса
	projectData, err := projectRepo.GetProjectByID(slotData.ProjectID)
	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}

	if projectData.OwnerID == userID {
		http.Error(w, "Нельзя откликнуться на собственный слот", http.StatusForbidden)
		return
	}

	if projectData.Status == "finished" {
		http.Error(w, "Проект завершён, заявки не принимаются", http.StatusConflict)
		return
	}

	// больше одной активной заявки нельзя
	checking, err := slotRepo.GetUserCheckingApplications(slotData.ID, userID)
	if err != nil {
		http.Error(w, "Ошибка проверки заявки", http.StatusInternalServerError)
		return
	}
	if checking != nil {
		http.Error(w, "Ранее заявка уже подавалась", http.StatusConflict)
		return
	}

	err = slotRepo.CreateApplication(slotData.ID, userID)
	if err != nil {
		http.Error(w, "Ошибка создания заявки", http.StatusInternalServerError)
		return
	}

	res := models.ApplicationResponse{
		Message: "Заявка подана",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
