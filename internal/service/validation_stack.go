package service

import (
	"HOTA/internal/repositories"
	"errors"
)

// все id стека должны быть в таблице stacks и не повторяться
func StacksExist(value any) error {
	stackIDs, ok := value.([]int)
	if !ok {
		return errors.New("стек должен быть списком id")
	}
	if len(stackIDs) == 0 {
		return nil // пустой список ловит validation.Required
	}

	return repositories.ValidateStacksExist(stackIDs)
}
