package project

import (
	repo "HOTA/internal/repositories/project"
	"encoding/json"
	"fmt"
	"net/http"
)

// отдается все публичные комнаты с статусом драфт и ворк
func GetProjectPublik(w http.ResponseWriter, r *http.Request) {

	rooms, err := repo.GetProjectPubliс()
	if err != nil {
		fmt.Printf("Ошибка получения комнат %v", err)
		http.Error(w, "Ошибка получения комнат %v", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)

}
