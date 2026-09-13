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
		http.Error(w, "Название комнаты обязателен", 500)
		return
	}
	if len(creat.Name) > 50 {
		http.Error(w, "Название комнаты не должно быть больше 50 символов", 500)
		return
	}

	if creat.Target == "" {
		http.Error(w, "Описание комнаты обязателен", 500)
		return
	}

	if len(creat.Target) > 300 {
		http.Error(w, "Описание комнаты не должно быть больше 300 символов", 500)
		return
	}

	if len(creat.Slots) > 6 {
		http.Error(w, "Максимум допустимо 6 слотов", 500)
		return
	}
	if len(creat.Slots) == 0 {
		http.Error(w, "Хотя бы 1 слот должен быть", 500)
		return
	}
	if creat.Privacy == "" {
		http.Error(w, "Поле privacy обязательно", http.StatusBadRequest)
		return
	}

	//комната
	roomModel := &models.Project{
		Name:        creat.Name,
		Description: creat.Target,
		OwnerID:     userID,
		Privacy:     creat.Privacy,
		Status:      "draft",
	}

	// создаем комнату
	if err = project.CreatProject(roomModel); err != nil {
		fmt.Printf("Ошибка создания комнаты %v\n", err)
		http.Error(w, "Ошибка создания комнаты", http.StatusInternalServerError)
		return
	}

	fmt.Println(roomModel.ID)

	// создать слоты и связь между комнатами
	var creatSlot []models.Slot

	for _, slot := range creat.Slots {
		slotModels := &models.Slot{
			ProjectID:  roomModel.ID,
			Rolle:   slot.Rolle,
			StackID: slot.Stack,
			Status:  "open",
		}

		if err = project.CreateSlot(slotModels); err != nil {
			fmt.Printf("Ошибка создания слота %v\n", err)
			http.Error(w, "Ошибка создания слота", http.StatusInternalServerError)
			return
		}
		creatSlot = append(creatSlot, *slotModels)
	}

	//ответ
	res := models.ProjectSlotOut{
		Project: *roomModel,
		Slots: creatSlot,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}
