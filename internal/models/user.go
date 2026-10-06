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
	ID       int     `json:"id"`
	Email    string  `json:"email,omitempty"`
	Nickname string  `json:"nickname"`
	Name     string  `json:"name"`
	Rolle    string  `json:"rolle"`
	Grade    string  `json:"grade"`
	Stack    []Stack `json:"stack"`
	About    string  `json:"about"`
}

type Stack struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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

type StackUser struct {
	StackID int `json:"stackID"`
}

type UpdateStacks struct {
	StackID []int `json:"stackID"`
}
