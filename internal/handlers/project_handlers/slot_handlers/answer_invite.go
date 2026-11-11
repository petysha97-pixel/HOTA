package slot

import (
	"HOTA/internal/models"
	slotRepo "HOTA/internal/repositories/project/slot"
	"HOTA/internal/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

// Кандидат принимает приглашение и попадает в очередь.
func AcceptInvite(w http.ResponseWriter, r *http.Request) {
	answerInvite(w, r, true)
}

// Кандидат отказывается от приглашения.
func DeclineInvite(w http.ResponseWriter, r *http.Request) {
	answerInvite(w, r, false)
}

// Проверяем пользователя и сохраняем его ответ.
func answerInvite(w http.ResponseWriter, r *http.Request, accept bool) {
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}
	inviteID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || inviteID <= 0 {
		http.Error(w, "Неверный ID приглашения", http.StatusBadRequest)
		return
	}

	// Добавление в очередь происходит только после принятия приглашения.
	err = slotRepo.AnswerInvite(inviteID, userID, accept)
	switch {
	case errors.Is(err, slotRepo.ErrInviteNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case errors.Is(err, slotRepo.ErrInviteForbidden):
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	case errors.Is(err, slotRepo.ErrInviteAnswered), errors.Is(err, slotRepo.ErrInviteUnavailable):
		http.Error(w, err.Error(), http.StatusConflict)
		return
	case err != nil:
		log.Printf("Ошибка ответа на приглашение: %v", err)
		http.Error(w, "Ошибка ответа на приглашение", http.StatusInternalServerError)
		return
	}

	message := "Приглашение отклонено"
	if accept {
		message = "Приглашение принято, вы добавлены в очередь"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ApplicationResponse{Message: message})
}
