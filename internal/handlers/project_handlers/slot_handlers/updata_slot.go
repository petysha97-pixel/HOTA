package slot

import (
	"HOTA/internal/models"
	slotRepo "HOTA/internal/repositories/project/slot"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// PATCH /slot/{id} — правка названия, описания, роли и стека, пока слот открыт и никто не утверждён
func UpdateSlot(w http.ResponseWriter, r *http.Request) {
	slotData, _, ok := ownerSlot(w, r)
	if !ok {
		return
	}

	var body models.DTOSlot
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := service.ValidateSlotStruct(body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if slotData.Status != "open" || slotData.UserID != nil {
		http.Error(w, "Менять можно только открытый слот без разработчика", http.StatusConflict)
		return
	}

	hasApproved, err := slotRepo.ApprovedAplecation(slotData.ID)
	if err != nil {
		http.Error(w, "Ошибка проверки заявки", http.StatusInternalServerError)
		return
	}
	if hasApproved {
		http.Error(w, "В слоте утверждён разработчик, менять нельзя", http.StatusConflict)
		return
	}

	if err = slotRepo.UpdateSlotInfo(slotData.ID, body.Name, body.Description, body.Rolle, body.StackID); err != nil {
		fmt.Printf("Ошибка обновления слота: %v\n", err)
		http.Error(w, "Ошибка обновления слота", http.StatusInternalServerError)
		return
	}

	updated, err := slotRepo.GetSlotByID(slotData.ID)
	if err != nil {
		http.Error(w, "Ошибка получения слота", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}
