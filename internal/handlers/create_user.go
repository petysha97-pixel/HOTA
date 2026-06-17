package handlers

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"bufio"
	"fmt"
	"os"
	"strings"
)

// регистрация нашего юзера через терминал
func CreateUser() {
	user := models.User{}

	fmt.Print("Ник: ")
	fmt.Scan(&user.Nickname)

	fmt.Print("Роль: ")
	fmt.Scan(&user.Role)

	var stack string
	fmt.Print("Стек (через запятую): ")
	fmt.Scan(&stack)
	user.Stack = strings.Split(stack, ",")

	fmt.Print("GitHub: ")
	fmt.Scan(&user.GitHub)

	fmt.Print("Telegram: ")
	fmt.Scan(&user.Telegram)

	fmt.Print("Статус: ")
	fmt.Scanln() // очищаем остаток после предыдущего Scan
	reader := bufio.NewReader(os.Stdin)
	status, _ := reader.ReadString('\n')
	user.Status = strings.TrimSpace(status)

	repositories.AppendUser(user)
	fmt.Println("Пользователь создан")
	fmt.Println()
}

// StatUser статистика пользователей
func StatUser() {
	count := repositories.CountUser()
	fmt.Printf("Количество зарегистрнированных пользователей: %d\n", count)

	stastRoll := repositories.CountByRole()
	fmt.Println("Статистика по ролям:")

	for role, count := range stastRoll {
		fmt.Printf("%s: %d\n", role, count)
	}
}

// ListUser Вывод всех пользователей
func ListUser() {
	users := repositories.GetAllUsers()
	if len(users) == 0 {
		fmt.Println("Зарегистрированных пользователей нет")
	} else {

		for _, user := range users {
			// strings.Join(user.Staсk, ",") //Преобразовываем из слайса в строку
			fmt.Printf("ID: %d,\nНик: %s,\nРоль: %s,\nСтек: %s, \nГитхаб: %s, \nTG: %s, \nСтатус: %s\n\n", user.ID, user.Nickname, user.Role, strings.Join(user.Stack, ","), user.GitHub, user.Telegram, user.Status)
		}

	}
}
