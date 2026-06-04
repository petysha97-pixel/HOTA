package main

import (
	"NOTA/internal/models"
	"NOTA/internal/repositories"
	"fmt"
)

func main() {
	fmt.Println(repositories.Users)
	for {
		fmt.Println("Введите действие")

		fmt.Println("1. Добавить пользователя")
		fmt.Println("2. Поиск по нику")
		fmt.Println("3. Поиск по id")
		fmt.Println("4. ​Вывести всех пользователей")
		fmt.Println("5. Обновление данных")

		var action int
		fmt.Scan(&action)
		fmt.Println()

		switch action {
		case 1:
			user := models.User{}

			fmt.Print("Ник: ")
			fmt.Scan(&user.Nickname)

			fmt.Print("Роль: ")
			fmt.Scan(&user.Role)

			// fmt.Print("Стек: ")
			// fmt.Scan(&user.Staсk)

			repositories.AppendUser(user) //nextID = 1
			fmt.Println("Пользователь создан")
			fmt.Println()
		case 2:
			var nickname string
			fmt.Println("Введите Ник: ")
			fmt.Scan(&nickname)

			user := repositories.SersheNickname(nickname)
			if user == nil {
				fmt.Println("Такого пользователя не существует")
				continue
			}
			fmt.Println(*user)

		case 3:

			var id int
			fmt.Println("Ведите ID")
			fmt.Scan(&id)

			user := repositories.SersheID(id)

			if user == nil {
				fmt.Printf("Такого пользователя с id: %d не существует", id)
				continue
			}
			fmt.Println(*user)

		case 4:
			users := repositories.GetAllUsers()
			for _, user := range users {
				fmt.Printf("ID: %d,\nНик: %s,\nРоль: %s\n", user.ID, user.Nickname, user.Role)
			}

		case 5:

			var id int
			fmt.Println("Введите ID")
			fmt.Scan(&id)

			user := repositories.SersheID(id)
			if user == nil {
				fmt.Printf("Такого пользователя с id: %d не существует", id)
				continue
			}
			// РЕДАКТИРОВАНИЯ ПОЛЬЗОВАТЕЛЯ

			var nik string
			fmt.Printf("Введите имя [%s]\n", user.Nickname)
			fmt.Scan(&nik)

			var rol string
			fmt.Printf("Введите новую роль [%s]\n", user.Role)
			fmt.Scan(&rol)

			user.Nickname = nik
			user.Role = rol

			ok := repositories.UpdateUser(*user)

			if ok == false {
				fmt.Println("Данного пользователя не удалось обновить")
			} else {
				fmt.Println("Пользовтель успешно обновлен")
			}

		default:
			fmt.Printf("Данного выбора: %d не существует", action)
		}

	}

}
