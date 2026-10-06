package project

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories/project"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

func CreatProject(w http.ResponseWriter, r *http.Request) {

	//берём их контекса + валидация айдишника
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	//запрос
	var creat models.DTOProject

	//парсим
	if err := json.NewDecoder(r.Body).Decode(&creat); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//валидация
	if creat.Name == "" {
		http.Error(w, "Название комнаты обязателен", http.StatusBadRequest)
		return
	}
	if len(creat.Name) > 50 {
		http.Error(w, "Название комнаты не должно быть больше 50 символов", http.StatusBadRequest)
		return
	}

	if creat.Description == "" {
		http.Error(w, "Описание комнаты обязателен", http.StatusBadRequest)
		return
	}

	if len(creat.Description) > 300 {
		http.Error(w, "Описание комнаты не должно быть больше 300 символов", http.StatusBadRequest)
		return
	}

	if len(creat.Slots) > 6 {
		http.Error(w, "Максимум допустимо 6 слотов", http.StatusBadRequest)
		return
	}
	if len(creat.Slots) == 0 {
		http.Error(w, "Хотя бы 1 слот должен быть", http.StatusBadRequest)
		return
	}
	if creat.Privacy == "" {
		http.Error(w, "Поле privacy обязательно", http.StatusBadRequest)
		return
	}

	// каждый слот: название, роль из каталога, стек
	for _, slot := range creat.Slots {
		if err := service.ValidateSlotStruct(slot); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	//комната
	roomModel := &models.Project{
		Name:        creat.Name,
		Description: creat.Description,
		OwnerID:     userID,
		Privacy:     creat.Privacy,
		Status:      "draft",
	}

	// слоты проекта: новый слот всегда открыт
	creatSlot := make([]models.Slot, 0, len(creat.Slots))
	for _, slot := range creat.Slots {
		creatSlot = append(creatSlot, models.Slot{
			Name:        slot.Name,
			Description: slot.Description,
			Rolle:       slot.Rolle,
			StackID:     slot.StackID,
			Status:      "open",
		})
	}

	// проект и слоты сохраняются одной транзакцией: при ошибке не остаётся ничего
	if err = project.CreateProjectWithSlots(roomModel, creatSlot); err != nil {
		fmt.Printf("Ошибка создания проекта %v\n", err)
		http.Error(w, "Ошибка создания проекта", http.StatusInternalServerError)
		return
	}

	//ответ
	res := models.ProjectSlotOut{
		Project: *roomModel,
		Slots:   creatSlot,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}
