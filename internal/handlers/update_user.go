package handlers

import (
	"HOTA/internal/repositories"
	"errors"
	"fmt"
)

func UpdateUser() error {
	var id int
	fmt.Println("Введите ID")
	fmt.Scan(&id)

	user := repositories.FindUserByID(id)
	if user == nil {
		return fmt.Errorf("Такого пользователя с id: %d не существует", id)
	}

	var nickname string
	fmt.Printf("Введите имя [%s]\n", user.Nickname)
	fmt.Scan(&nickname)

	var role string
	fmt.Printf("Введите новую роль [%s]\n", user.Role)
	fmt.Scan(&role)

	var gh string
	fmt.Printf("Введите GitHub [%s]\n", user.GitHub)
	fmt.Scan(&gh)

	var tg string
	fmt.Printf("Введите Telegram [%s]\n", user.Telegram)
	fmt.Scan(&tg)

	var status string
	fmt.Printf("Введите статус [%s]\n", user.Status)
	fmt.Scan(&status)

	user.Nickname = nickname
	user.Role = role
	user.GitHub = gh
	user.Telegram = tg
	user.Status = status

	ok := repositories.UpdateUser(*user)

	if ok == false {
		return errors.New("Данного пользователя не удалось обновить")
	}

	return nil
}
