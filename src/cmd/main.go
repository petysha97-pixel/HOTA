package main

import (
	"HOTA/src/internal/constants"
	"HOTA/src/internal/handlers/getters"
	"HOTA/src/internal/handlers/setters"
	"HOTA/src/internal/models"
	"HOTA/src/internal/utils"
)

var menu models.Menu = models.Menu{
	Entries: []models.MenuEntry{
		{
			Name: constants.Green + "Создать пользователя" + constants.Reset,
			Func: setters.CreateUser,
		},
		{
			Name: constants.Cyan + "Обновить пользователя" + constants.Reset,
			Func: setters.UpdateUser,
		},
		{
			Name: constants.Red + "Удалить пользователя" + constants.Reset,
			Func: setters.DeleteUser,
		},
		{
			Name: constants.Yellow + "Поиск по ID" + constants.Reset,
			Func: getters.GetUserById,
		},
		{
			Name: constants.Yellow + "Поиск по нику" + constants.Reset,
			Func: getters.GetUserByNickname,
		},
		{
			Name: constants.Yellow + "Поиск по стэку" + constants.Reset,
			Func: getters.GetUserByStack,
		},
		{
			Name: constants.Blue + "Список пользователей" + constants.Reset,
			Func: getters.ListUser,
		},
		{
			Name: constants.Blue + "Статистика" + constants.Reset,
			Func: getters.StatUser,
		},
	},
}

var firstIteration bool = true

func main() {
	for {
		if(firstIteration) {
			utils.PrintSeparator(0, 1, 50, "=")
		} else {
			utils.PrintSeparator(1, 1, 50, "=")
		}

		utils.PrintMenu(menu)
	
		entry := utils.ReadPositiveInt("Номер: ", true, false)
		
		if entry > 0 && entry <= len(menu.Entries) {
			utils.PrintSeparator(1, 1, 50, "=")
			menu.Entries[entry - 1].Func()
		}

		if(firstIteration == true) {
			firstIteration = false
		}
	}
}
