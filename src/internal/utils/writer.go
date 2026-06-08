package utils

import (
	"HOTA/src/internal/constants"
	"HOTA/src/internal/models"
	"fmt"
	"strconv"
)

func PrintSeparator(indentsUp, indentsDown, lenght int, symbol string) {
	for i := 0; i < indentsUp; i++ {
		fmt.Println()
	}
	
	for i := 0; i < lenght; i++ {
		fmt.Print(symbol)
	}

	for i := 0; i <= indentsDown; i++ {
		fmt.Println()
	}
}

func PrintValue(property, value string) {
	fmt.Printf(constants.Cyan + "%s: " + constants.Green + "%s\n" + constants.Reset, property, value)
}

func PrintMultipleValue(property string, values []string) {
	fmt.Printf("%s:\n", property)

	for _, value := range values {
		fmt.Printf(constants.Cyan + "  - %s: " + constants.Green + "%s\n" + constants.Reset, property, value)
	}
}

func PrintMenu(menu models.Menu) {
	for key, value := range menu.Entries {
		fmt.Printf("%d. %s\n", key + 1, value.Name)
	}
}

func PrintUser(user models.User) {
	PrintValue("ID", strconv.Itoa(user.Id))
	PrintValue("Nickname", user.Nickname)
	PrintValue("Role", user.Role)
	PrintMultipleValue("Stack", user.Stack)
	PrintValue("Github", user.GitHub)
	PrintValue("Telegram", user.Telegram)
	PrintValue("Status", user.Status)
}
