package slot

import (
	"HOTA/internal/models"
	slotRepo "HOTA/internal/repositories/project/slot"
	"HOTA/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

// Показываем разработчику его входящие.
func GetInbox(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Получаем приглашения и уведомления текущего пользователя.
	response, err := slotRepo.GetInbox(userID)
	if err != nil {
		log.Printf("Ошибка входящих: %v", err)
		http.Error(w, "Ошибка получения входящих", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Помечаем собственные уведомления прочитанными.
func ReadNotifications(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}
	if err = slotRepo.ReadNotifications(userID); err != nil {
		log.Printf("Ошибка прочтения уведомлений: %v", err)
		http.Error(w, "Ошибка прочтения уведомлений", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ApplicationResponse{Message: "Уведомления прочитаны"})
}
