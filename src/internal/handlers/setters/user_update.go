package setters

import (
	"HOTA/src/internal/constants"
	"HOTA/src/internal/repositories"
	"HOTA/src/internal/utils"
	"fmt"
	"strings"
)

func UpdateUser() {
	user := repositories.GetById(utils.ReadPositiveInt("ID: ", false, true))
	
	if user == nil {
		fmt.Println(constants.Red + "Пользователь не найден!" + constants.Reset)
	} else {
		user.Nickname = utils.ReadString(fmt.Sprintf("Ник [%s]: ", user.Nickname), false, false)
		user.Role = utils.ReadString(fmt.Sprintf("Роль [%s]: ", user.Role), false, false)
		user.Stack = strings.Split(utils.ReadString(fmt.Sprintf("Стэк (через запятую, без пробелов) [%s]: ", strings.Join(user.Stack, ",")), false, false), ",")
		user.GitHub = utils.ReadString(fmt.Sprintf("GitHub [%s]: ", user.GitHub), false, false)
		user.Telegram = utils.ReadString(fmt.Sprintf("Telegram [%s]: ", user.Telegram), false, false)
		user.Status = utils.ReadString(fmt.Sprintf("Статус [%s]: ", user.Status), false, true)
		
		if repositories.Update(*user) {
			fmt.Println(constants.Green + "Пользователь обновлён!" + constants.Reset)
		} else {
			fmt.Println(constants.Red + "Не удалось обновить пользователя!" + constants.Reset)
		}
	}
}
