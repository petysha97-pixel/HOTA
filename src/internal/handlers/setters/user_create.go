package setters

import (
	"HOTA/src/internal/constants"
	"HOTA/src/internal/models"
	"HOTA/src/internal/repositories"
	"HOTA/src/internal/utils"
	"fmt"
	"strings"
)

func CreateUser() {
	user := models.User{
		Nickname: utils.ReadString("Ник: ", false, false),
		Role: utils.ReadString("Роль: ", false, false),
		Stack: strings.Split(utils.ReadString("Стэк (через запятую, без пробелов): ", false, false), ","),
		GitHub: utils.ReadString("GitHub: ", false, false),
		Telegram: utils.ReadString("Telegram: ", false, false),
		Status: utils.ReadString("Статус: ", false, true),
	}

	repositories.Add(user)

	fmt.Println(constants.Green + "Пользователь создан!" + constants.Reset)
}