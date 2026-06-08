package getters

import (
	"HOTA/src/internal/constants"
	"HOTA/src/internal/repositories"
	"HOTA/src/internal/utils"
	"fmt"
)

func ListUser() {
	users := repositories.Users
	
	if len(users) == 0 {
		fmt.Println(constants.Red + "Пока пользователей нет!" + constants.Reset)
	} else {
		for key, user := range users {
			if key != 0 {
				utils.PrintSeparator(1, 1, 15, "-")
			}

			utils.PrintUser(user)
		}
	}
}