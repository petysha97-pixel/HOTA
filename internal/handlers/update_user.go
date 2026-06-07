package handlers

import (
	"HOTA/internal/repositories"
	"fmt"
)

func UpdateUser() {
	var id int
	fmt.Println("Введите ID")
	fmt.Scan(&id)

	user := repositories.SersheID(id)
	if user == nil {
		fmt.Printf("Пользователь с id: %d не найден\n", id)
		return
	}

	var nik string
	fmt.Printf("Введите имя [%s]\n", user.Nickname)
	fmt.Scan(&nik)

	var rol string
	fmt.Printf("Введите новую роль [%s]\n", user.Role)
	fmt.Scan(&rol)

	var gh string
	fmt.Printf("Введите GitHub [%s]\n", user.GitHub)
	fmt.Scan(&gh)

	var tg string
	fmt.Printf("Введите Telegram [%s]\n", user.Telegram)
	fmt.Scan(&tg)

	var stats string
	fmt.Printf("Введите статус [%s]\n", user.Status)
	fmt.Scan(&stats)

	user.Nickname = nik
	user.Role = rol
	user.GitHub = gh
	user.Telegram = tg
	user.Status = stats

	ok := repositories.UpdateUser(*user)

	if ok == false {
		fmt.Printf("Пользователь с id: %d не найден\n", id)
		return
	}

	fmt.Println("Пользователь успешно обновлен")
}
