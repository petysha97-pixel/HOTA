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

// тело запроса: только новый статус
type slotStatusIn struct {
	Status string `json:"status"`
}

// обновляем статус слота
func UpdateStatusSlot(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// ID слота из URL: /slot/{id}/status
	slotID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Неверный ID слота", http.StatusBadRequest)
		return
	}

	var body slotStatusIn
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if body.Status != "open" && body.Status != "close" && body.Status != "done" {
		http.Error(w, "Неизвестный статус слота", http.StatusBadRequest)
		return
	}

	slotData, err := slotRepo.GetSlotByID(slotID)
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
	if slotData.Status == body.Status {
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
		if body.Status != "close" {
			http.Error(w, "Из open можно только в close", http.StatusConflict)
			return
		}
		if !hasApproved {
			http.Error(w, "Нельзя закрыть слот: нет принятой заявки", http.StatusConflict)
			return
		}
	case "review":
		// Возвращаем кандидата в очередь и открываем слот.
		if body.Status != "open" {
			http.Error(w, "Из review можно вернуться в open", http.StatusConflict)
			return
		}
	case "close":
		if body.Status == "open" {
			// переоткрытие всегда, принятую заявку снимаем
		} else if body.Status == "done" {
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

	// Сохраняем новый статус. При переоткрытии снимаем исполнителя.
	err = slotRepo.UpdateSlotByID(slotData.ID, userID, slotData.Status, body.Status)
	if errors.Is(err, slotRepo.ErrSlotStatusConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Ошибка обновления статуса", http.StatusInternalServerError)
		return
	}

	// ответ
	res := models.UpdateStatusRoom{
		Message: "Статус обновлён",
		Status:  body.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
