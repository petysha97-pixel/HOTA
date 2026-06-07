package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"fmt"
	"strings"
)

func CreateUser() {
	user := models.User{}

	fmt.Print("Ник: ")
	fmt.Scan(&user.Nickname)

	fmt.Print("Роль: ")
	fmt.Scan(&user.Role)

	var stack string
	fmt.Print("Стек (через запятую): ")
	fmt.Scan(&stack)
	user.Staсk = strings.Split(stack, ",")

	fmt.Print("GitHub: ")
	fmt.Scan(&user.GitHub)

	fmt.Print("Telegram: ")
	fmt.Scan(&user.Telegram)

	fmt.Print("Статус: ")
	fmt.Scan(&user.Status)

	repositories.CreateUser(user)

	fmt.Println("Пользователь создан")
	fmt.Println()
}

// Statistics статистика пользователей
func Statistics() {
	count := repositories.CountUser()
	fmt.Printf("Количество зарегистрнированных пользователей: %d\n", count)

	roleStat := repositories.CountByRole()
	fmt.Println("Статистика по ролям:")

	for role, count := range roleStat {
		fmt.Printf("%s: %d\n", role, count)
	}
}

// ListUsers Вывод всех пользователей
func ListUsers() {
	users := repositories.GetAllUsers()
	if len(users) == 0 {
		fmt.Println("Зарегистрированных пользователей нет")
		return
	}

	for _, user := range users {
		// strings.Join(user.Staсk, ",") //Преобразовываем из слайса в строку
		fmt.Printf(
			"ID: %d,\nНик: %s,\nРоль: %s,\nСтек: %s, \nГитхаб: %s, \nTG: %s, \nСтатус: %s",
			user.ID,
			user.Nickname,
			user.Role,
			strings.Join(user.Staсk, ","),
			user.GitHub,
			user.Telegram,
			user.Status,
		)
	}
}
