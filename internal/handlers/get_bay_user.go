package handlers

import (
	"HOTA/internal/printer"
	"HOTA/internal/repositories"
	"fmt"
	"strings"
)

// GetUserByNik поиск по нику
func GetUserByNik() {
	var nickname string
	fmt.Println("Введите Ник: ")
	fmt.Scan(&nickname)

	user := repositories.SersheNickname(nickname)
	if user == nil {
		fmt.Println("Такого пользователя не существует")
		return
	}
	printer.PrintUser(*user)
}

// GetUserByIDПоиск по ID
func GetUserByID() {
	var id int
	fmt.Println("Ведите ID")
	fmt.Scan(&id)

	user := repositories.SersheID(id)

	if user == nil {
		fmt.Printf("Такого пользователя с id: %d не существует", id)
		return
	}
	printer.PrintUser(*user)
}

// SerheUserByStack Поиск по стеку
func SerheUserByStack() {
	var stack string
	fmt.Println("Введите стек: ")
	fmt.Scan(&stack)

	users := repositories.FindByStack(stack)

	if len(users) == 0 {
		fmt.Println("Ничего не найдено")
	}
	for _, user := range users {
		fmt.Printf("ID: %d,\nНик: %s,\nРоль: %s,\nСтек: %s\n", user.ID, user.Nickname, user.Role, strings.Join(user.Stack, ","))

	}
}

func DeleteIserID() {
	var id int
	fmt.Println("Введите ID")
	fmt.Scan(&id)

	ok := repositories.DeleteUser(id)

	if ok == false {
		fmt.Printf("Такого пользователя с id: %d не существует", id)
	} else {
		fmt.Println("Пользователь успешно удален")
	}
}
