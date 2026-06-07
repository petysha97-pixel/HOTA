package main

import (
	"HOTA/internal/handlers"
	"HOTA/internal/printer"
	"fmt"
)

func Menu() int {
	fmt.Println()
	printer.PrintMeny()
	fmt.Println()

	var action int
	fmt.Scan(&action)
	fmt.Println()
	return action
}

func main() {

	actions := []func(){
		func() { handlers.CreateUser() },
		func() { handlers.GetUserByNik() },
		func() { handlers.GetUserByID() },
		func() { handlers.SerheUserByStack() },
		func() {
			err := handlers.UpdateUser()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println("Пользователь успешно обновлен")
		},
		func() { handlers.DeleteIserID() },
		func() { handlers.StatUser() },
	}

	for {
		action := Menu() + 1

		if action > 0 && action <= len(actions) {
			actions[action-1]()
		} else {

			fmt.Printf("Данного выбора: %d не существует", action)
		}
	}

}
