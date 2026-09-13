package project

import (
	"HOTA/internal/repositories/project"
	"encoding/json"
	"fmt"
	"net/http"
)

// отдается все публичные комнаты с статусом драфт и ворк
func GetProjectPubliс(w http.ResponseWriter, r *http.Request) {

	rooms, err := project.GetProjectPublik()
	if err != nil {
		fmt.Printf("Ошибка получения комнат %v", err)
		http.Error(w, "Ошибка получения комнат %v", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rooms)

}
