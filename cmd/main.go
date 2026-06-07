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
	actions := map[int]func(){
		1: handlers.CreateUser,
		2: handlers.GetUserByNik,
		3: handlers.GetUserByID,
		4: handlers.SerheUserByStack,
		5: handlers.ListUser,
		6: handlers.UpdateUser,
		7: handlers.DeleteIserID,
		8: handlers.StatUser,
	}

	for {
		action := Menu()

		if action < 1 || action > len(actions) {
			fmt.Println("Данного действия нет, попробуйте еще раз")
			continue
		}

		actions[action]()

	}

}
