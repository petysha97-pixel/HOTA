package slot

import (
	"HOTA/internal/models"
	slotRepo "HOTA/internal/repositories/project/slot"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// Отправляем разработчику приглашение на слот.
func InviteDeveloper(w http.ResponseWriter, r *http.Request) {
	// Проверяем слот и права владельца проекта.
	slotData, projectData, ok := ownerSlot(w, r)
	if !ok {
		return
	}

	// Получаем ID приглашённого разработчика.
	applicantID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil || applicantID <= 0 {
		http.Error(w, "Неверный ID разработчика", http.StatusBadRequest)
		return
	}

	// Создаём ожидающее приглашение.
	err = slotRepo.InviteToSlot(slotData.ID, applicantID, projectData.OwnerID)
	if errors.Is(err, slotRepo.ErrInviteConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Ошибка приглашения: %v", err)
		http.Error(w, "Ошибка отправки приглашения", http.StatusInternalServerError)
		return
	}

	// Сообщаем фронту об отправке приглашения.
	response := models.ApplicationResponse{Message: "Приглашение отправлено"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
