package getters

import (
	"HOTA/src/internal/repositories"
	"fmt"
)

func StatUser() {
	fmt.Printf("Кол-во пользователей: %d\n", len(repositories.Users))

	for role, count := range repositories.CountByRole() {
		fmt.Printf("  - %s: %d\n", role, count)
	}
}