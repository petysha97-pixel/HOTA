package handlers

import (
	"HOTA/internal/printer"
	"HOTA/internal/repositories"
	"fmt"
	"strings"
)

// GetUserByNickname поиск по нику
func GetUserByNickname() {
	var nickname string
	fmt.Println("Введите Ник: ")
	fmt.Scan(&nickname)

	user := repositories.FindUserByNickname(nickname)
	if user == nil {
		fmt.Println("Такого пользователя не существует")
		return
	}
	printer.PrintUser(*user)
}

// GetUserByID Поиск по ID
func GetUserByID() {
	var id int
	fmt.Println("Введите ID")
	fmt.Scan(&id)

	user := repositories.FindUserByID(id)

	if user == nil {
		fmt.Printf("Такого пользователя с id: %d не существует", id)
		return
	}
	printer.PrintUser(*user)
}

// FindUsersByStack Поиск по стеку
func GetUsersByStack() {
	var stack string
	fmt.Println("Введите стек: ")
	fmt.Scan(&stack)

	users := repositories.FindByStack(stack)

	if len(users) == 0 {
		fmt.Println("Ничего не найдено")
	}
	for _, user := range users {
		fmt.Printf(
			"ID: %d,\nНик: %s,\nРоль: %s,\nСтек: %s\n",
			user.ID,
			user.Nickname,
			user.Role,
			strings.Join(user.Staсk, ","),
		)
	}
}

// DeleteUserByID удаление пользователя по ID
func DeleteUserByID() {
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
