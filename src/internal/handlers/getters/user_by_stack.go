package getters

import (
	"HOTA/src/internal/constants"
	"HOTA/src/internal/repositories"
	"HOTA/src/internal/utils"
	"fmt"
)

func GetUserByStack() {
	users := repositories.GetByStack(utils.ReadString("Стэк: ", false, true))

	if len(users) == 0 {
		fmt.Println(constants.Red + "Пользователь не найден!" + constants.Reset)
	} else {
		for key, user := range users {
			if key != 0 {
				utils.PrintSeparator(1, 1, 15, "-")
			}

			utils.PrintUser(user)
		}
	}
}