package models

import (
	"database/sql"
	"time"
)

// Стуктура нашего пользователя
type User struct {
	ID       int       `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Nickname string    `json:"nickname"`
	Name     string    `json:"name"` // ФИО, необязательное
	Rolle    string    `json:"rolle"`
	Grade    string    `json:"grade"` // Junior / Middle / Senior
	StackID  []int     `json:"stackID"`
	About    string    `json:"about"`
	CreatAt  time.Time `json:"creatAt"`
	UpdateAt time.Time `json:"updateAt"`
}

// DTO для ответа (без пароля)
type UserResponse struct {
	ID       int           `json:"id"`
	Email    string        `json:"email,omitempty"`
	Nickname string        `json:"nickname"`
	Name     string        `json:"name"`
	Rolle    string        `json:"rolle"`
	Grade    string        `json:"grade"`
	Stack    []Stack       `json:"stack"`
	About    string        `json:"about"`
	History  []HistoryItem `json:"history,omitempty"` // история участия, только в профиле
}

// технология пользователя: описание опыта он пишет в профиле
type Stack struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"` // описание опыта с технологией
}

// одна строка «Истории участия» в профиле
// type = "project" — проект, который пользователь создал
// type = "slot" — слот, где пользователь исполнитель
type HistoryItem struct {
	Type        string    `json:"type"`
	ProjectID   int       `json:"projectID"`
	ProjectName string    `json:"projectName"`
	Status      string    `json:"status"` // у проекта: draft/working/finished, у слота: open/review/close/done
	CreatAt     time.Time `json:"creatAt"`

	// только для проекта: сколько слотов всего, сколько занято (исполнитель есть, работа не сдана), сколько сдано
	SlotsTotal int `json:"slotsTotal"`
	SlotsTaken int `json:"slotsTaken"`
	SlotsDone  int `json:"slotsDone"`

	// только для слота
	SlotID        int    `json:"slotID,omitempty"`
	SlotName      string `json:"slotName,omitempty"`
	Rolle         string `json:"rolle,omitempty"`
	StackID       []int  `json:"stackID,omitempty"`
	OwnerNickname string `json:"ownerNickname,omitempty"` // кто создал проект (кто принял в слот)
}

// правка профиля из окна «Профиль»: почта и стек меняются своими ручками
type UpdateProfile struct {
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
	Rolle    string `json:"rolle"`
	Grade    string `json:"grade"`
	About    string `json:"about"`
}

// DTO структура для ответа стека при регистрации
type UserStackID struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Роль из каталога ролей
type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var UserDB *sql.DB

// DTO структура для ответа при обновлении пароля
type UpdatePasswors struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// DTO структура для обновления email
type UpdateEmail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// DTO структура авторизации ПОЛУЧЕНИЕ
type AuotIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// DTO структура авторизации ОТВЕТА
type AuotOut struct {
	Token string `json:"token"`
}

// описание о себе
type About struct {
	About string `json:"about"`
}

// добавить или удалить одну технологию в стеке пользователя
// description нужен только при добавлении, его можно не передавать
type StackUser struct {
	StackID     int    `json:"stackID"`
	Description string `json:"description"`
}

// изменить описание опыта у одной технологии
type StackChange struct {
	Description string `json:"description"`
}

type UpdateStacks struct {
	StackID []int `json:"stackID"`
}
