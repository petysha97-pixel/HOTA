package models

type User struct {
	Id       int
	Nickname string
	Role     string
	Stack    []string
	GitHub   string
	Telegram string
	Status   string
}
