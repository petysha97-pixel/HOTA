package main

import (
	"HOTA/internal/handlers"
	"HOTA/internal/printer"
	"fmt"
)

func Menu() int {
	fmt.Println()
	printer.PrintMenu()
	fmt.Println()

	var action int
	fmt.Scan(&action)
	fmt.Println()
	return action
}

func main() {

	for {
		action := Menu()

		switch action {
		case 1: // Добавить пользователя
			handlers.CreateUser()

		case 2: // Поиск по нику
			handlers.GetUserByNickname()

		case 3: // Поиск по ID
			handlers.GetUserByID()

		case 4: // Поиск по стеку
			handlers.GetUsersByStack()

		case 5: // Вывести всех пользователей
			handlers.ListUsers()

		case 6: // Обновление данных
			err := handlers.UpdateUser()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			fmt.Println("Пользователь успешно обновлен")
		case 7: // Удаление пользователя
			handlers.DeleteUserByID()

		case 8: //статистика
			handlers.Statistics()

		default:
			fmt.Printf("Данного выбора: %d не существует", action)
		}

	}

}
