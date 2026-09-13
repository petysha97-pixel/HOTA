package slot

import (
	projectRepo "HOTA/internal/repositories/project"
	slotRepo "HOTA/internal/repositories/project/slot"
	"HOTA/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

// получить все заявки на слот
func GetSlotApplications(w http.ResponseWriter, r *http.Request) {
	// проверяем пользователя
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// достаём ID слота из URL
	slotID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Неверный ID слота", http.StatusBadRequest)
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
		http.Error(w, "Только владелец проекта может смотреть заявки", http.StatusForbidden)
		return
	}

	// получаем все заявки на слот
	apps, err := slotRepo.GetApplicationsBySlot(slotID)
	if err != nil {
		http.Error(w, "Ошибка получения заявок", http.StatusInternalServerError)
		return
	}

	// отдаём список заявок
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apps)
}
