package models

import "time"

// фильтры вкладки «Проекты» (приходят из адресной строки)
type ProjectFilter struct {
	Filter string // all — все, mine — мои (создал или в команде), role — есть открытый слот моей роли
	Status string // draft / working / finished, пусто — все кроме finished
	Query  string // текст: название, описание, ник автора, роли слотов
	Rolle  string // моя роль — нужна для filter=role (берём из профиля)
}

// слот в карточке проекта
type SlotShort struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"` // техзадание
	Rolle       string  `json:"rolle"`
	Status      string  `json:"status"` // open / review / close / done
	UserID      *int    `json:"userID"` // исполнитель, null — никого
	Stack       []Stack `json:"stack"`
}

// сколько слотов в каждом статусе
type ProjectCounts struct {
	Open   int `json:"open"`
	Review int `json:"review"`
	Close  int `json:"close"`
	Done   int `json:"done"`
	Total  int `json:"total"`
}

// карточка проекта: вкладка «Проекты» и страница проекта
type ProjectCard struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Status        string        `json:"status"`  // draft / working / finished
	Privacy       string        `json:"privacy"` // public / private
	CreatAt       time.Time     `json:"creatAt"`
	OwnerID       int           `json:"ownerID"`
	OwnerNickname string        `json:"ownerNickname"`
	Slots         []SlotShort   `json:"slots"`
	Counts        ProjectCounts `json:"counts"`
}

// слот в полоске левой панели
type RailSlot struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	UserID *int   `json:"userID"`
}

// проект в левой панели «Проекты»
type RailProject struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Status     string     `json:"status"`
	Slots      []RailSlot `json:"slots"`
	SlotsTotal int        `json:"slotsTotal"` // всего слотов
	SlotsTaken int        `json:"slotsTaken"` // занято: есть исполнитель, работа не сдана
	MySlots    []string   `json:"mySlots"`    // названия моих слотов (для «Участвую»)

	// цифра-бейдж и из чего она сложена
	Applications int `json:"applications"` // отклики ждут ответа (только у создателя)
	Review       int `json:"review"`       // слоты на знакомстве (только у создателя)
	Invites      int `json:"invites"`      // приглашения мне ждут ответа
	Badge        int `json:"badge"`        // сумма трёх цифр выше
}

// ответ левой панели
type MyProjectsResponse struct {
	Created     []RailProject `json:"created"`     // «Создал»
	Participate []RailProject `json:"participate"` // «Участвую»
}
