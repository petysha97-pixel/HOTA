package utils

import (
	"HOTA/src/internal/constants"
	"fmt"
)

func ReadString(prompt string, indentUp, indentDown bool) string {
	var input string

	if indentUp {
		fmt.Println()
	}

	fmt.Printf(constants.Purple + "%s" + constants.Yellow, prompt)
	fmt.Scan(&input)
	fmt.Print(constants.Reset)

	if indentDown {
		fmt.Println()
	}

	return input;
}

func ReadPositiveInt(prompt string, indentUp, indentDown bool) int {
	input := -1

	if indentUp {
		fmt.Println()
	}

	for input < 0 {
		fmt.Printf(constants.Purple + "%s" + constants.Yellow, prompt)
		fmt.Scan(&input)
		fmt.Print(constants.Reset)
	}

	if indentDown {
		fmt.Println()
	}

	return input;
}