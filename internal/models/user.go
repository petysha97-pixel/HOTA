package models

// import "encoding/json"

// Стуктура нашего пользователя
type User struct {
	ID       int      `json:"ID"`
	Nickname string   `json:"nickname"`
	Role     string   `json:"role"`
	Staсk    []string `json:"staсk"`
	GitHub   string   `json:"github"`   //ссылка гитхаб json
	Telegram string   `json:"telegram"` //ссылка тг json
	Status   string   `json:"status"`   //в сети или нет json

}
