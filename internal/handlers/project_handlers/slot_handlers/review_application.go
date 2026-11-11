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

// Берём кандидата из очереди на знакомство.
func ReviewApplication(w http.ResponseWriter, r *http.Request) {
	// Проверяем слот, владельца проекта и доступ к изменениям.
	slotData, projectData, ok := ownerSlot(w, r)
	if !ok {
		return
	}

	// Получаем ID кандидата из адреса запроса.
	applicantID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil || applicantID <= 0 {
		http.Error(w, "Неверный ID кандидата", http.StatusBadRequest)
		return
	}

	// Назначаем кандидата и переводим слот в review.
	err = slotRepo.ReviewSlotApplication(
		slotData.ID,
		applicantID,
		projectData.OwnerID,
	)
	if errors.Is(err, slotRepo.ErrReviewConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("Ошибка знакомства: %v", err)
		http.Error(w, "Ошибка назначения кандидата", http.StatusInternalServerError)
		return
	}

	// Сообщаем фронту, что кандидат назначен.
	response := models.ApplicationResponse{
		Message: "Кандидат взят на знакомство",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
