package handlers

import (
	"HOTA/internal/repositories"
	"encoding/json"
	"fmt"
	"net/http"
)

// берем все роли для выбора при регистрации, в профиле и в слоте
func GetRoles(w http.ResponseWriter, r *http.Request) {

	roles, err := repositories.GetRoles()
	if err != nil {
		fmt.Printf("Ошибка получения ролей: %v\n", err)
		http.Error(w, "Ошибка получения ролей", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}
