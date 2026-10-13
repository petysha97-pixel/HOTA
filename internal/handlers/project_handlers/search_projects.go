package project

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	projectRepo "HOTA/internal/repositories/project"
	"HOTA/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"unicode/utf8"
)

// вкладка «Проекты»: список и поиск
// GET /projects?filter=all&status=draft&query=шлюз
// все параметры необязательные
func SearchProjects(w http.ResponseWriter, r *http.Request) {

	userID, err := service.ContextUserIDValid(r)
	if err != nil {
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	// собираем фильтры из адресной строки
	filter := models.ProjectFilter{
		Filter: r.URL.Query().Get("filter"),
		Status: r.URL.Query().Get("status"),
		Query:  strings.TrimSpace(r.URL.Query().Get("query")),
	}
	if filter.Filter == "" {
		filter.Filter = "all" // по умолчанию — «Все»
	}

	// проверяем параметры
	if filter.Filter != "all" && filter.Filter != "mine" && filter.Filter != "role" {
		http.Error(w, "filter: all, mine или role", http.StatusBadRequest)
		return
	}
	if filter.Status != "" && filter.Status != "draft" && filter.Status != "working" && filter.Status != "finished" {
		http.Error(w, "status: draft, working или finished", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(filter.Query) > 50 {
		http.Error(w, "текст поиска не должен быть больше 50 символов", http.StatusBadRequest)
		return
	}

	// мой профиль: роль нужна для «Моя роль» и для сортировки
	me := repositories.Get_userdb(userID)
	filter.Rolle = me.Rolle

	// мой стек — для сортировки «сначала то, что мне подходит»
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

	projects, err := projectRepo.SearchProjects(filter, userID)
	if err != nil {
		fmt.Printf("Ошибка получения проектов %v\n", err)
		http.Error(w, "Ошибка получения проектов", http.StatusInternalServerError)
		return
	}

	// текстовый поиск: название, описание, ник автора, роли слотов
	result := []models.ProjectCard{}
	for _, p := range projects {
		if projectTextFound(filter.Query, p) {
			result = append(result, p)
		}
	}

	// сначала проекты, которые мне больше подходят
	sort.SliceStable(result, func(i, j int) bool {
		return podhodit(result[i], me.Rolle, uMenya) > podhodit(result[j], me.Rolle, uMenya)
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// есть ли текст поиска в названии, описании, нике автора или ролях слотов
// ищем без учёта больших и маленьких букв (работает и для русских)
func projectTextFound(text string, p models.ProjectCard) bool {
	if text == "" {
		return true // текста нет — подходит любой
	}

	vse := []string{p.Name, p.Description, p.OwnerNickname}
	for _, slot := range p.Slots {
		vse = append(vse, slot.Rolle)
	}

	stroka := strings.ToLower(strings.Join(vse, " "))
	return strings.Contains(stroka, strings.ToLower(text))
}

// насколько проект мне подходит: смотрим только открытые слоты
// +1 за каждую мою технологию в слоте, +2 если роль слота — моя
func podhodit(p models.ProjectCard, myRolle string, uMenya map[int]bool) int {
	ochki := 0
	for _, slot := range p.Slots {
		if slot.Status != "open" {
			continue
		}
		for _, s := range slot.Stack {
			if uMenya[s.ID] {
				ochki++
			}
		}
		if slot.Rolle == myRolle {
			ochki += 2
		}
	}
	return ochki
}
