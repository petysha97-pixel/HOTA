package models


// Стуктура нашего пользователя
type User struct {
	ID       int
	Nickname string
	Role     string
	Staсk    []string
}