package slot

import (
	"HOTA/internal/models"
	slotRepo "HOTA/internal/repositories/project/slot"
	"encoding/json"
	"fmt"
	"strconv"

	"HOTA/internal/service"
	"net/http"
)

// создание заявки на слот
func ApplicationsSlot(w http.ResponseWriter, r *http.Request) {
	// берём из контекста + валидация
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	slotID := r.PathValue("id")
	fmt.Println(slotID)
	slotIDint, err := strconv.Atoi(slotID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Некорректный id слота", http.StatusBadRequest)
		fmt.Printf("Не получилось конвертировать строку в число: %s", err)
		return
	}

	//проверка слота в БД
	slot, err := slotRepo.GetSlotByID(slotIDint)
	if err != nil {
		http.Error(w, "Слот не найден", http.StatusUnauthorized)
		return
	}

	if slot.Status != "open" {
		http.Error(w, "Слот закрыт", http.StatusUnauthorized)
		return
	}

	// проверка заявки на слот (пользователь не может пожавать более чем 1 заявку)
	CheckingApplications, _ := slotRepo.GetUserCheckingApplications(slot.ID, userID)
	if CheckingApplications != nil {
		fmt.Printf("дубль подачи заявки на слот: %s", err)
		http.Error(w, "Ранее заявка уже подавалась", 409)
	}

	//подача заявки
	err = slotRepo.CreateApplication(slot.ID, userID)
	if err != nil {
		http.Error(w, "Ошибка создания заявки", http.StatusInternalServerError)
		return
	}

	res := models.ApplicationResponse{
		Message: "Заявка подана",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
