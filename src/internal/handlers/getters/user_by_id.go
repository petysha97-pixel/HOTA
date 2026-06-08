package getters

import (
	"HOTA/src/internal/constants"
	"HOTA/src/internal/repositories"
	"HOTA/src/internal/utils"
	"fmt"
)

func GetUserById() {
	user := repositories.GetById(utils.ReadPositiveInt("ID: ", false, true))

	if user == nil {
		fmt.Println(constants.Red + "Пользователь не найден!" + constants.Reset)
	} else {
		utils.PrintUser(*user)
	}
}