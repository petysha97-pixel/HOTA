package slot

import (
	"HOTA/internal/models"
	slotRepo "HOTA/internal/repositories/project/slot"
	"encoding/json"
	"fmt"
	"net/http"
)

// DELETE /slot/{id} — удалить слот вместе с его стеком и заявками
func DeleteSlot(w http.ResponseWriter, r *http.Request) {
	slotData, _, ok := ownerSlot(w, r)
	if !ok {
		return
	}

	if err := slotRepo.DeleteSlotByID(slotData.ID); err != nil {
		fmt.Printf("Ошибка удаления слота: %v\n", err)
		http.Error(w, "Ошибка удаления слота", http.StatusInternalServerError)
		return
	}

	res := models.ApplicationResponse{
		Message: "Слот удалён",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
