package main

import (
	"HOTA/internal/handlers"
	"HOTA/internal/printer"
	"HOTA/internal/repositories"
	"fmt"

	"strings"
)

func main() {
	
	// Массив функций
	var actions = []func(){
		func() {
			handlers.CreateUser()
			fmt.Println("Пользователь создан")
			fmt.Println()
		},
		func() {
			var nickname string
			fmt.Println("Введите Ник: ")
			fmt.Scan(&nickname)

			user := repositories.SersheNickname(nickname)
			if user == nil {
				fmt.Println("Такого пользователя не существует")
				return
			}
			printer.PrintUser(*user)
		},
		func() {
			var id int
			fmt.Println("Ведите ID")
			fmt.Scan(&id)

			user := repositories.SersheID(id)
			if user == nil {
				fmt.Printf("Такого пользователя с id: %d не существует\n", id)
				return
			}
			printer.PrintUser(*user)
		},
		func() {
			var stack string
			fmt.Println("Введите стек: ")
			fmt.Scan(&stack)

			users := repositories.FindByStack(stack)
			if len(users) == 0 {
				fmt.Println("Ничего не найдено")
			}
			for _, user := range users {
				fmt.Printf("ID: %d,\nНик: %s,\nРоль: %s,\nСтек: %s\n", user.ID, user.Nickname, user.Role, strings.Join(user.Staсk, ","))
			}
		},
		func() {
			users := repositories.GetAllUsers()
			for _, user := range users {
				fmt.Printf("ID: %d,\nНик: %s,\nРоль: %s,\nСтек: %s, \nГитхаб: %s, \nTG: %s, \nСтатус: %s\n", 
					user.ID, user.Nickname, user.Role, strings.Join(user.Staсk, ","), user.GitHub, user.Telegram, user.Status)
			}
		},
		func() {
			err := handlers.UpdateUser()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("Пользователь успешно обновлен")
		},
		func() {
			var id int
			fmt.Println("Введите ID")
			fmt.Scan(&id)

			ok := repositories.DeleteUser(id)
			if !ok {
				fmt.Printf("Такого пользователя с id: %d не существует\n", id)
			} else {
				fmt.Println("Пользователь успешно удален")
			}
		},
		func() {
			count := repositories.CountUser()
			fmt.Printf("Количество зарегистрированных пользователей: %d\n", count)

			stastRoll := repositories.CountByRole()
			fmt.Println("Статистика по ролям:")
			for role, count := range stastRoll {
				fmt.Printf("%s: %d\n", role, count)
			}
		},
	}

	for {
		fmt.Println()
		printer.PrintMeny()
		fmt.Println()

		var action int
		fmt.Scan(&action)
		fmt.Println()

		// Проверка номера экшена
		if action >= 0 && action < len(actions) {
			actions[action]()
		} else {
			fmt.Printf("Данного выбора: %d не существует", action)
		}
	}
}
