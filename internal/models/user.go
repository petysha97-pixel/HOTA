package models

import (
	"database/sql"
	"time"
)

// Стуктура нашего пользователя
type User struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Nickname   string `json:"nickname"`
	Rolle      string `json:"rolle"`
	Stack      []int  `json:"stack"`
	About      string `json:"about"`
	Creat_add  time.Time
	Update_add time.Time
}

// DTO для ответа (без пароля)
type UserResponse struct {
	ID       int     `json:"id"`
	Email    string  `json:"email,omitempty"`
	Nickname string  `json:"nickname"`
	Rolle    string  `json:"rolle"`
	Stack    []Stack `json:"stack"`
	About    string  `json:"about"`
}

type Stack struct {
	ID   int
	Name string
}

// DTO структура для ответа стека при регистрации
type UserStackID struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var UserDB *sql.DB

// DTO структура для ответа при обновлении пароля
type UpdatePasswors struct {
	OldPassword string `json:"oldpassword"`
	NewPassword string `json:"newpassword"`
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

//описание о себе
type Description struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}



type StackUser struct {
	IDStack int `json:"id_stack"`
}

type UpdateStacks struct {
    StackID []int `json:"stack_id"`
}


