package project

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// драфт - воркинг - финиш

// обновляем статус комнаты
func UpdateStatusProject(w http.ResponseWriter, r *http.Request) {

	// берём из контекста + валидация
	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// Достем id из гет-запроса
	projectID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Не верный ID проекта", http.StatusBadRequest)
		return
	}

	//парсим
	var bodyProject models.Project
	if err := json.NewDecoder(r.Body).Decode(&bodyProject); err != nil {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	// достаём комнату из БД по ID из тела
	projectData, err := projectRepo.GetProjectByID(projectID) //какой айди использователь? из гет запроса или из тела паршеного
	if err != nil {
		http.Error(w, "Комната не найдена", http.StatusNotFound)
		return
	}

	//установка статуса только создателем
	if userID != projectData.OwnerID {
		http.Error(w, "смена статуса запрещена", 403)
		return

	}

	//запрет на установления одного и того же статуса
	if bodyProject.Status == projectData.Status {
		http.Error(w, "Данный статус уже установлен", 403)
		return
	}

	// посчитать общее колиство слотов и по отдельности сколько открытых и сколько закрытых
	totalSlot, openSlot, _, doneSlot, err := projectRepo.CountSlotStatus(projectID)
	if err != nil {
		fmt.Printf("Ошибка подсчета слотов %v", err)
		http.Error(w, "Ошибка подсчета слотов", http.StatusNotFound)
		return
	}

	switch projectData.Status {
	case "draft":
		if bodyProject.Status != "working" {
			http.Error(w, "из draft можно только в working ", 403)
			return
		}

		if openSlot > 0 {
			http.Error(w, "перевод статуса запрещен, есть открытые слоты ", 403)
			return
		}

	case "working":
		if bodyProject.Status == "draft" {
			//создатель открыл набор
		} else if bodyProject.Status == "finished" {
			if doneSlot != totalSlot {
				http.Error(w, "перевод статуса запрещен, не все слоты завершены", 403)
				return
			}
		}

	case "finished":
		http.Error(w, "Проект завершен, запрещена смена статуса", 403)
		return

	default:
		if bodyProject.Status != "draft" && bodyProject.Status != "working" && bodyProject.Status != "finished" {
			http.Error(w, "Неизвестный статус проекта", 403)
			return
		}

	}

	err = projectRepo.UpdateProjectStatus(projectID, bodyProject.Status)
	if err != nil {
		http.Error(w, "ошибка обновления стататуса проекта", 403)
		return
	}

	// ответ
	req := models.UpdateStatusRoom{
		Message: "Статус обновлён",
		Status:  bodyProject.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(req)
}

//anykey
//nota-app.dev
