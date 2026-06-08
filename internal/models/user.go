package models

// Стуктура нашего пользователя
type User struct {
	ID       int `json:"id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Staсk    []string `json:"stack"`
	GitHub   string `json:"githubLink"` // ссылка гитхаб
	Telegram string `json:"telegramLink"` // ссылка тг
	Status   string `json:"status"` // в сети или нет
}

type Users struct {
	Count int `json:"count"`
	Data []User `json:"data"`
}