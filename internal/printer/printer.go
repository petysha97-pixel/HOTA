package printer

import (
	"HOTA/internal/models"
	"fmt"
)

// PrintMeny - вывод меню в терминал
func PrintMeny() {
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

// PrintUser - вывод информацию о пользователе в терминал
func PrintUser(user models.User) {
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
	fmt.Printf("Статус: %s\n", user.Status)

	fmt.Println("​-----------------------------------------")

}
