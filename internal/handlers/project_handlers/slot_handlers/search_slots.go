package slot

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	slotRepo "HOTA/internal/repositories/project/slot"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// вкладка «Слоты»: слоты, куда ещё можно попасть
// GET /slots?tab=open&query=api&rolle=Backend&stackID=1,5&onlyFree=true&fullMatch=true
// все параметры необязательные
func SearchSlots(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// собираем фильтры из адресной строки
	filter := models.SlotFilter{
		Tab:       r.URL.Query().Get("tab"),
		Query:     strings.TrimSpace(r.URL.Query().Get("query")),
		Rolle:     r.URL.Query().Get("rolle"),
		OnlyFree:  r.URL.Query().Get("onlyFree") == "true",
		FullMatch: r.URL.Query().Get("fullMatch") == "true",
	}
	if filter.Tab == "" {
		filter.Tab = "open" // по умолчанию — «Открытые роли»
	}

	// stackID приходит строкой "1,5,7" — переводим в список чисел
	stackText := r.URL.Query().Get("stackID")
	if stackText != "" {
		for _, part := range strings.Split(stackText, ",") {
			id, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				http.Error(w, "stackID — это id технологий через запятую", http.StatusBadRequest)
				return
			}
			filter.StackID = append(filter.StackID, id)
		}
	}

	if err := service.ValidateSlotFilter(filter); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// мой стек — чтобы посчитать совпадение
	myStack, err := repositories.GetStacksByUserID(userID)
	if err != nil {
		fmt.Printf("Ошибка получения моего стека %v\n", err)
		http.Error(w, "Ошибка получения стека", http.StatusInternalServerError)
		return
	}
	uMenya := map[int]bool{} // id технологий, которые у меня есть
	for _, s := range myStack {
		uMenya[s.ID] = true
	}

	// слоты, где мой отклик ждёт ответа
	applied, err := slotRepo.GetMyPendingSlotIDs(userID)
	if err != nil {
		fmt.Printf("Ошибка получения откликов %v\n", err)
		http.Error(w, "Ошибка получения откликов", http.StatusInternalServerError)
		return
	}

	slots, err := slotRepo.SearchSlots(filter, userID)
	if err != nil {
		fmt.Printf("Ошибка получения слотов %v\n", err)
		http.Error(w, "Ошибка получения слотов", http.StatusInternalServerError)
		return
	}

	response := models.SlotSearchResponse{Slots: []models.SlotCard{}}

	for _, slot := range slots {

		// текстовый поиск: роль, название, техзадание, технологии
		if !slotTextFound(filter.Query, slot) {
			continue
		}

		// сравниваем стек слота с моим
		slot.StackMissing = []string{}
		for _, s := range slot.Stack {
			if uMenya[s.ID] {
				slot.StackHits++
			} else {
				slot.StackMissing = append(slot.StackMissing, s.Name)
			}
		}

		// полное совпадение: у меня есть все технологии слота
		if filter.FullMatch && (len(slot.Stack) == 0 || len(slot.StackMissing) > 0) {
			continue
		}

		slot.Applied = applied[slot.ID]

		// «Открытые роли»: слот свободен или на знакомстве с кем-то другим (не со мной)
		yaKandidat := slot.UserID != nil && *slot.UserID == userID
		otkryt := !yaKandidat

		// цифры на вкладках считаем с теми же фильтрами
		if otkryt {
			response.OpenCount++
		}
		if slot.Applied {
			response.MineCount++
		}

		// оставляем только слоты выбранной вкладки
		if filter.Tab == "open" && !otkryt {
			continue
		}
		if filter.Tab == "mine" && !slot.Applied {
			continue
		}

		response.Slots = append(response.Slots, slot)
	}

	// сначала слоты, где больше совпадение с моим стеком
	sort.SliceStable(response.Slots, func(i, j int) bool {
		a, b := response.Slots[i], response.Slots[j]
		dolyaA, dolyaB := slotDolya(a), slotDolya(b)
		if dolyaA != dolyaB {
			return dolyaA > dolyaB
		}
		return a.StackHits > b.StackHits
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// есть ли текст поиска в роли, названии, техзадании или технологиях слота
// ищем без учёта больших и маленьких букв (работает и для русских)
func slotTextFound(text string, slot models.SlotCard) bool {
	if text == "" {
		return true // текста нет — подходит любой
	}

	vse := []string{slot.Rolle, slot.Name, slot.Description}
	for _, s := range slot.Stack {
		vse = append(vse, s.Name)
	}

	stroka := strings.ToLower(strings.Join(vse, " "))
	return strings.Contains(stroka, strings.ToLower(text))
}

// доля совпадения стека: 0 — ничего, 1 — всё
func slotDolya(slot models.SlotCard) float64 {
	if len(slot.Stack) == 0 {
		return 0
	}
	return float64(slot.StackHits) / float64(len(slot.Stack))
}
