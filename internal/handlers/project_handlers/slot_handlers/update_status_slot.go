package slot

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	slotRepo "HOTA/internal/repositories/project/slot"
	"HOTA/internal/service"
	"encoding/json"
	"net/http"
)

// обновляем статус слота
func UpdateStatusSlot(w http.ResponseWriter, r *http.Request) {
	// берём из контекста + валидация
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// парсим тело
	var slot models.Slot
	if err := json.NewDecoder(r.Body).Decode(&slot); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// проверяем, существет ли данный проект вообще
	projectDa, err := projectRepo.GetProjectByID(slot.ProjectID)
	if slot.ProjectID != projectDa.ID {
		http.Error(w, "данного проекта не существует", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Комната не найдена", http.StatusNotFound)
		return
	}

	NewSlot, err := slotRepo.GetSlotByID(slot.ID)
	if err != nil {
		http.Error(w, "слот не найден", 500)
		return
	}

	project, err := projectRepo.GetProjectByID(NewSlot.ProjectID)
	if err != nil {
		http.Error(w, "ошибка получение проекта у данного слота", 500)
		return
	}

	if userID != project.OwnerID {
		http.Error(w, "столько Создатель имеет право менять статус слота", 500)
		return
	}

	if project.Status == "finished" {
		http.Error(w, "смена статуса не доступна, когда проект завершен", 500)
		return
	}

	// все статусы слота
	if NewSlot.Status != "open" && NewSlot.Status != "close" && NewSlot.Status != "done" {
		http.Error(w, "Допустимые статусы: open, close, done", 500)
		return
	}

	switch NewSlot.Status {
	case "open":
		if slot.Status == "done" {
			http.Error(w, "Нельзя перейти из open в done, сначала закройте слот (closed)", 500)
			return
		}
	case "done":
		if slot.Status != "done" {
			http.Error(w, "Нельзя изменить статус выполненного слота", http.StatusBadRequest)
			return
		}
	}

	_, err = slotRepo.UpdateSlotByID(NewSlot.ID, slot.Status)
	if err != nil {
		http.Error(w, "ошибка обновления статуса", 500)
		return
	}

	// ответ
	req := models.UpdateStatusRoom{
		Message: "Статус обновлён",
		Status:  NewSlot.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)

}
