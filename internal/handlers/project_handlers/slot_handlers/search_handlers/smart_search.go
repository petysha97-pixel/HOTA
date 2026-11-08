package slot

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	projectRepo "HOTA/internal/repositories/project"
	slotRepo "HOTA/internal/repositories/project/slot"
	candRepo "HOTA/internal/repositories/project/slot/searhe_repo"
	"HOTA/internal/service"
)

func SearchCandidates(w http.ResponseWriter, r *http.Request) {

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

	// проверка слота в БД
	slotData, err := slotRepo.GetSlotByID(slotID)
	if err != nil {
		http.Error(w, "Слот не найден", http.StatusNotFound)
		return
	}

	// проект слота: нужен для проверки владельца и статуса
	projectData, err := projectRepo.GetProjectByID(slotData.ProjectID)
	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}

	if projectData.OwnerID != userID {
		http.Error(w, "Только владелец проекта может искать разработчиков", http.StatusForbidden)
		return
	}

	if projectData.Status == "finished" {
		http.Error(w, "Проект закрыт, прислать кандидатов нельзя", http.StatusConflict)
		return
	}

	if slotData.Status != "open" {
		http.Error(w, "Слот уже закрыт либо занят", http.StatusConflict)
		return
	}

	slotData.StackID, err = projectRepo.GetStackIDsBySlotID(slotData.ID)
	if err != nil {
		http.Error(w, "Ошибка при получении стека слота", http.StatusInternalServerError)
		return
	}

	// стеки: по умолчанию стеки слота, пользователь мог снять чипы
	stackIDs := slotData.StackID
	if raw := r.URL.Query().Get("stack"); raw != "" {
		slotStacks := make(map[int]bool, len(slotData.StackID))
		for _, id := range slotData.StackID {
			slotStacks[id] = true
		}

		stackIDs = nil
		for _, part := range strings.Split(raw, ",") {
			id, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				http.Error(w, "Неверный список стеков", http.StatusBadRequest)
				return
			}
			// берём только стеки, которые есть у слота
			if slotStacks[id] {
				stackIDs = append(stackIDs, id)
			}
		}
	}

	// фильтры: ник или имя и грейд
	query := Normalize(strings.TrimSpace(r.URL.Query().Get("q")))
	grade := strings.ToLower(r.URL.Query().Get("grade"))

	candidates, err := candRepo.FindCandidates(slotData.Rolle, stackIDs, userID, query, grade)
	if err != nil {
		http.Error(w, "Ошибка при подборе кандидатов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(candidates)
}

func Normalize(text string) string {
	return strings.ReplaceAll(strings.ToLower(text), "ё", "е")
}
