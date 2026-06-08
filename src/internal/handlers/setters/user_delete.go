package setters

import (
	"HOTA/src/internal/constants"
	"HOTA/src/internal/repositories"
	"HOTA/src/internal/utils"
	"fmt"
)

func DeleteUser() {
	if repositories.Delete(utils.ReadPositiveInt("ID: ", false, true)) {
		fmt.Println(constants.Green + "Пользователь удалён!" + constants.Reset)
	} else {
		fmt.Println(constants.Red + "Не удалось удалить пользователя!" + constants.Reset)
	}
}