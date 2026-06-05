package main

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"fmt"

	"strings"
)

func main() {

	for {
		fmt.Println("​-----------------------------------------")
		fmt.Println()
		fmt.Println("Введите действие:")
		fmt.Println("1. Добавить пользователя")
		fmt.Println("2. Поиск по нику")
		fmt.Println("3. Поиск по id")
		fmt.Println("4. Поиск по стеку")
		fmt.Println("5. ​Вывести всех пользователей")
		fmt.Println("6. Обновление данных")
		fmt.Println("7. Удалить пользователя")
		fmt.Println("8. Статистика")
		fmt.Println()
		fmt.Println("​-----------------------------------------")

		var action int
		fmt.Scan(&action)
		fmt.Println()

		switch action {
		case 1: // Добавить пользователя
			user := models.User{}

			fmt.Print("Ник: ")
			fmt.Scan(&user.Nickname)

			fmt.Print("Роль: ")
			fmt.Scan(&user.Role)

			var stack string
			fmt.Print("Стек (через запятую): ")
			fmt.Scan(&stack)

			user.Staсk = strings.Split(stack, ",")

			repositories.AppendUser(user)
			fmt.Println("Пользователь создан")
			fmt.Println()
		case 2: // Поиск по нику
			var nickname string
			fmt.Println("Введите Ник: ")
			fmt.Scan(&nickname)

			user := repositories.SersheNickname(nickname)
			if user == nil {
				fmt.Println("Такого пользователя не существует")
				continue
			}
			fmt.Println(*user)

		case 3: // Поиск по ID

			var id int
			fmt.Println("Ведите ID")
			fmt.Scan(&id)

			user := repositories.SersheID(id)

			if user == nil {
				fmt.Printf("Такого пользователя с id: %d не существует", id)
				continue
			}
			fmt.Println(*user)
		case 4: // Поиск по стеку
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

		case 5: // Вывести всех пользователей
			users := repositories.GetAllUsers()
			for _, user := range users {
				// strings.Join(user.Staсk, ",") //Преобразовываем из слайса в строку
				fmt.Printf("ID: %d,\nНик: %s,\nРоль: %s,\nСтек: %s\n", user.ID, user.Nickname, user.Role, strings.Join(user.Staсk, ","))
			}

		case 6: // Обновление данных

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

		case 7: // Удаление пользователя
			var id int
			fmt.Println("Введите ID")
			fmt.Scan(&id)

			ok := repositories.DeleteUser(id)

			if ok == false {
				fmt.Printf("Такого пользователя с id: %d не существует", id)
			} else {
				fmt.Println("Пользователь успешно удален")
			}
		case 8: //статистика
			count := repositories.CountUser()
			fmt.Printf("Количество зарегистрнированных пользователей: %d\n", count)

			stastRoll := repositories.CountByRole()
			fmt.Println("Статистика по ролям:")

			for role, count := range stastRoll {
				fmt.Printf("%s: %d\n", role, count)
			}

		default:
			fmt.Printf("Данного выбора: %d не существует", action)
		}

	}

}
