package slot

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	slotRepo "HOTA/internal/repositories/project/slot"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// снимаем разработчика со слота при ревью
func RemoveSlotMember(w http.ResponseWriter, r *http.Request) {

	// проверяем пользователя
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// ID слота из URL: /slot/{id}/status
	slotID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || slotID <= 0 {
		http.Error(w, "Неверный ID слота", http.StatusBadRequest)
		return
	}

	// проверяем, что слот существует и возвращаем заполненные слот обратно
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

	//доступ к снятию только у владельца проекта
	if projectData.OwnerID != userID {
		http.Error(w, "Только владелец проекта может снимать разработчика со слота ", http.StatusForbidden)
		return
	}

	//снимаем разработчика со слота только в статусе open
	if slotData.Status != "open" && slotData.Status != "review" {
		http.Error(w, "слот закрыт, разработчика нельзя снять", http.StatusConflict)
		return
	}

	//снимаем разработчика со слота
	count, err := slotRepo.RemoveApprovedMember(slotData.ID)
	if err != nil {
		fmt.Printf("ошибка снятия разработчика со слота %v", err)
		http.Error(w, "Некого удалять со слота, он пустой", http.StatusConflict)
		return
	}

	if count == 0 {
		fmt.Printf("Слот и так пустой, некого удалять %v", err)
		http.Error(w, "Слот и так пустой, некого удалять", http.StatusConflict)
		return
	}


	// ответ
	res := models.ApplicationResponse{
		Message: "Разработчик снят со слота",
	}

	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}
