package service

import (
	"HOTA/internal/models"
	"HOTA/internal/repositories"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

// проверка описания опыта у технологии: по желанию, до 300 символов
func ValidateStackInfo(info models.StackChange) error {
	return validation.ValidateStruct(&info,
		validation.Field(&info.Description, validation.Length(0, 300).Error("описание опыта не должно быть больше 300 символов")),
	)
}
