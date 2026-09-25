package handlers

import (
	"HOTA/internal/repositories"
	"encoding/json"
	"net/http"
)

//берем все стеки для выбора при регистрации
func GetStacks(w http.ResponseWriter, r *http.Request) {

	stacks, err := repositories.GetStacks()
	if err != nil {
		http.Error(w, "ошибка получения стека в регистрации", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stacks)

}
