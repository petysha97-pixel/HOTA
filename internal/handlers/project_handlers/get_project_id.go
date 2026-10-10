package project

import (
	"HOTA/internal/models"
	roomRepo "HOTA/internal/repositories/project"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// отдаём комнату со слотами (любой статус)
func GetProject(w http.ResponseWriter, r *http.Request) {
	// вытаскиваем ID из пути /project/123
	GetIDstring := r.PathValue("id")
	fmt.Println(GetIDstring)
	roomID, err := strconv.Atoi(GetIDstring)
	if err != nil {
		http.Error(w, "Неверный ID комнаты", http.StatusBadRequest)
		return
	}

	room, err := roomRepo.GetProjectByID(roomID)
	if err != nil {
		http.Error(w, "Комната не найдена", http.StatusNotFound)
		return
	}

	if room.Privacy == "private" {
		userID, err := service.ContextUserIDValid(r)
		if err != nil {
			http.Error(w, "Доступ запрещён", http.StatusForbidden)
			return
		}
		if room.OwnerID != userID {
			http.Error(w, "Доступ запрещён", http.StatusForbidden)
			return
		}
	}

	slots, err := roomRepo.GetSlotsByProjectID(roomID)
	if err != nil {
		fmt.Printf("Ошибка в поиске слота %v\n", err)
		http.Error(w, "Ошибка получения слотов", http.StatusInternalServerError)
		return
	}

	//ответ
	res := models.ProjectSlotOut{
		Project: *room,
		Slots:   slots,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
