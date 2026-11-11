package slot

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	slotRepo "HOTA/internal/repositories/project/slot"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"

	"HOTA/internal/service"
	"net/http"
)

// тело запроса на отклик: сообщение создателю проекта (можно не писать)
type applyIn struct {
	Message string `json:"message"`
}

// создание заявки на слот
func ApplicationsSlot(w http.ResponseWriter, r *http.Request) {
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	slotID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Некорректный id слота", http.StatusBadRequest)
		return
	}

	// проверка слота в БД
	slotData, err := slotRepo.GetSlotByID(slotID)
	if err != nil {
		http.Error(w, "Слот не найден", http.StatusNotFound)
		return
	}

	// Отклики принимаем, пока ищем человека или знакомимся с кандидатом.
	if slotData.Status != "open" && slotData.Status != "review" {
		http.Error(w, "Набор на слот закрыт", http.StatusConflict)
		return
	}

	// проект слота: нужен для проверки владельца и статуса
	projectData, err := projectRepo.GetProjectByID(slotData.ProjectID)
	if err != nil {
		http.Error(w, "Проект не найден", http.StatusNotFound)
		return
	}

	if projectData.OwnerID == userID {
		http.Error(w, "Нельзя откликнуться на собственный слот", http.StatusForbidden)
		return
	}

	if projectData.Status == "finished" {
		http.Error(w, "Проект завершён, заявки не принимаются", http.StatusConflict)
		return
	}

	// сообщение к заявке по желанию: тело может быть пустым
	// io.EOF значит, что тело пустое — это нормально, сообщение не обязательно
	var body applyIn
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil && !errors.Is(err, io.EOF) {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if utf8.RuneCountInString(body.Message) > 500 {
		http.Error(w, "Сообщение не должно быть больше 500 символов", http.StatusBadRequest)
		return
	}

	// заявка на слот у пользователя одна (UNIQUE в БД): смотрим, что с прошлой
	prev, err := slotRepo.GetUserApplication(slotData.ID, userID)
	if err != nil {
		http.Error(w, "Ошибка проверки заявки", http.StatusInternalServerError)
		return
	}

	if prev == nil {
		// раньше не откликался — создаём новую заявку
		err = slotRepo.CreateApplication(slotData.ID, userID, body.Message)

	} else if prev.Status == "pending" {
		// заявка уже ждёт решения
		http.Error(w, "Ранее заявка уже подавалась", http.StatusConflict)
		return

	} else if prev.Status == "approved" {
		// уже утверждён на этот слот
		http.Error(w, "Вы уже утверждены на этот слот", http.StatusConflict)
		return

	} else {
		// Возвращаем отклонённую или снятую заявку в очередь.
		err = slotRepo.ReopenApplication(prev.ID, body.Message)
	}

	// Слот могли закрыть, пока обрабатывался отклик.
	if errors.Is(err, slotRepo.ErrApplicationConflict) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		fmt.Printf("Ошибка создания заявки: %v\n", err)
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
