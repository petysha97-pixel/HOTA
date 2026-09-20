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
	var bodySlot models.Slot
	if err := json.NewDecoder(r.Body).Decode(&bodySlot); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// достаём слот из БД
	slotData, err := slotRepo.GetSlotByID(bodySlot.ID)
	if err != nil {
		http.Error(w, "Слот не найден", http.StatusNotFound)
		return
	}

	// достаём проект этого слота
	projectData, err := projectRepo.GetProjectByID(slotData.ProjectID)
	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}

	// статус слота меняет только создатель проекта
	if userID != projectData.OwnerID {
		http.Error(w, "Только создатель может менять статус слота", http.StatusForbidden)
		return
	}

	// в завершённом проекте слоты не трогаем
	if projectData.Status == "finished" {
		http.Error(w, "Проект завершён, слоты менять нельзя", http.StatusConflict)
		return
	}

	// нельзя поставить тот же статус
	if slotData.Status == bodySlot.Status {
		http.Error(w, "Статус уже установлен", http.StatusConflict)
		return
	}

	// есть ли принятая заявка на слот
	hasApproved, err := slotRepo.ApprovedAplecation(slotData.ID)
	if err != nil {
		http.Error(w, "Ошибка проверки заявки", http.StatusInternalServerError)
		return
	}

	// проверяем переход по текущему статусу
	switch slotData.Status {
	case "open":
		// из open только в close и только при принятой заявке
		if bodySlot.Status != "close" {
			http.Error(w, "Из open можно только в close", http.StatusConflict)
			return
		}
		if !hasApproved {
			http.Error(w, "Нельзя закрыть слот: нет принятой заявки", http.StatusConflict)
			return
		}

	case "close":
		if bodySlot.Status == "open" {
			// переоткрытие всегда, принятую заявку снимаем
		} else if bodySlot.Status == "done" {
			if !hasApproved {
				http.Error(w, "Нельзя выполнить слот: нет принятой заявки", http.StatusConflict)
				return
			}
		} else {
			http.Error(w, "Из close можно в open или done", http.StatusConflict)
			return
		}

	case "done":
		http.Error(w, "Слот выполнен, статус не меняется", http.StatusConflict)
		return

	default:
		http.Error(w, "Неизвестный статус слота", http.StatusBadRequest)
		return
	}

	// при переоткрытии снимаем принятую заявку
	if slotData.Status == "close" && bodySlot.Status == "open" {
		err = slotRepo.ResetApprovedApplication(slotData.ID)
		if err != nil {
			http.Error(w, "Ошибка снятия заявки", http.StatusInternalServerError)
			return
		}
	}

	// обновляем статус в БД
	bodySlot2, err := slotRepo.UpdateSlotByID(slotData.ID, bodySlot.Status)
	if err != nil {
		http.Error(w, "Ошибка обновления статуса", http.StatusInternalServerError)
		return
	}

	// ответ
	res := models.UpdateStatusRoom{
		Message: "Статус обновлён",
		Status:  bodySlot2.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
