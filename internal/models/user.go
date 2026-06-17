package models

// import "encoding/json"

// Стуктура нашего пользователя
type User struct {
	ID       int      `json:"id"`
	Nickname string   `json:"nickname"`
	Role     string   `json:"role"`
	Stack    []string `json:"stack"`
	GitHub   string   `json:"github"`   //ссылка гитхаб json
	Telegram string   `json:"telegram"` //ссылка тг json
	Status   string   `json:"status"`   //в сети или нет json

}


//возвращаемая струкрута с ошибками
type Registration_fail struct {
	StatusGlobal string `json:"statusglobal"`
	Error []string  `json:"error"`

}