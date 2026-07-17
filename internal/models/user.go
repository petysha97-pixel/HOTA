package models

import (
	"database/sql"
	"time"
)

// Стуктура нашего пользователя
type User struct {
	ID         int      `json:"id"`
	Login      string   `json:"login"`
	Password   string   `json:"password"`
	Nickname   string   `json:"nickname"`
	Rolle      string   `json:"rolle"`
	Stack      []string `json:"stack"`
	Creat_add  time.Time
	Update_add time.Time
}


// DTO для ответа (без пароля)
type UserResponse struct {
	ID       int      `json:"id"`
	Nickname string   `json:"nickname"`
	Rolle    string   `json:"rolle"`
	Stack    []string `json:"stack"`
}


var UserDB *sql.DB
