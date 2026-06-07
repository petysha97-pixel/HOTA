package services

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"fmt"

	"strings"
)

func Run() {
	for {
		fmt.Println()
		printMeny()
        fmt.Println()
		var action int
		fmt.Scan(&action)
		fmt.Println()

		switch action {
		case 1: // Добавить пользователя
			newUser()

		case 2: // Поиск по нику
			printUserByNickname()

		case 3: // Поиск по ID
			printUserByID()

		case 4: // Поиск по стеку
			printUsersByStack()

		case 5: // Вывести всех пользователей
			printAllUsers()

		case 6: // Обновление данных
			updateUserByID()

		case 7: // Удаление пользователя
			deleteUserByID()

		case 8: //статистика
			printStats()

		default:
			fmt.Printf("Данного выбора: %d не существует", action)
		}
	}
}

// фукнция вызова меню
func printMeny() {
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
}

//фукнция вывода пользователя
func printUser(user models.User) {
	fmt.Println("​-----------------------------------------")
	fmt.Printf("ID: %d\n", user.ID)
	fmt.Printf("Ник: %s\n", user.Nickname)
	fmt.Printf("Роль: %s\n", user.Role)

	fmt.Println("Стек:")
	for _, stek := range user.Stack {
		fmt.Printf("%s\n", stek)
	}

	fmt.Printf("GitHub: %s\n", user.GitHub)
	fmt.Printf("Telegram: %s\n", user.Telegram)
	fmt.Printf("Статус: %s\n", user.Telegram)

	fmt.Println("​-----------------------------------------")
}

// newUser создает нового пользователя из ввода терминала
func newUser() {
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
	fmt.Scan(&user.Status)
	repositories.AppendUser(user)

	fmt.Println("Пользователь создан")
	fmt.Println()
}

// printUserByNickname выводит пользователя по нику из ввода терминала
func printUserByNickname() {
	var nickname string
	fmt.Println("Введите Ник: ")
	fmt.Scan(&nickname)

	user := repositories.SersheNickname(nickname) // тут только один пользователь может вернуться?
	if user == nil {
		fmt.Println("Такого пользователя не существует")
		return
	}
	printUser(*user)
}

// printUserByID выводит пользователя по ID из ввода терминала
func printUserByID() {
	var id int
	fmt.Println("Ведите ID")
	fmt.Scan(&id)

	user := repositories.SersheID(id)

	if user == nil {
		fmt.Printf("Такого пользователя с id: %d не существует", id)
		return
	}
	printUser(*user)
}

// printUsersByStack выводит пользователей по стеку из ввода терминала
func printUsersByStack() {
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

// printAllUsers выводит всех пользователей
func printAllUsers() {
	users := repositories.GetAllUsers()
	for _, user := range users {
		// strings.Join(user.Staсk, ",") //Преобразовываем из слайса в строку
		fmt.Printf("ID: %d,\nНик: %s,\nРоль: %s,\nСтек: %s, \nГитхаб: %s, \nTG: %s, \nСтатус: %s", user.ID, user.Nickname, user.Role, strings.Join(user.Stack, ","), user.GitHub, user.Telegram, user.Status)
	}
}

// updateUserByID обновляет пользователя по ID из ввода терминала
func updateUserByID() {
	var id int
	fmt.Println("Введите ID")
	fmt.Scan(&id)

	user := repositories.SersheID(id)
	if user == nil {
		fmt.Printf("Такого пользователя с id: %d не существует", id)
		return
	}
	// РЕДАКТИРОВАНИЯ ПОЛЬЗОВАТЕЛЯ

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
		fmt.Println("Данного пользователя не удалось обновить")
	} else {
		fmt.Println("Пользовтель успешно обновлен")
	}
}

// deleteUserByID удаляет пользователя по ID из ввода терминала
func deleteUserByID() {
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

// printStats выводит статистику по всем пользователям
func printStats() {
	count := repositories.CountUser()
	fmt.Printf("Количество зарегистрнированных пользователей: %d\n", count)

	stastRoll := repositories.CountByRole()
	fmt.Println("Статистика по ролям:")

	for role, count := range stastRoll {
		fmt.Printf("%s: %d\n", role, count)
	}
}