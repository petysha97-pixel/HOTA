package service

import (
	"HOTA/internal/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// проверка фильтров поиска разработчиков
func ValidateUserFilter(filter models.UserFilter) error {
	return validation.ValidateStruct(&filter,

		// текст поиска
		validation.Field(&filter.Query, validation.Length(0, 50).Error("текст поиска не должен быть больше 50 символов")),

		// роль — из каталога ролей
		validation.Field(&filter.Rolle, validation.By(RoleExists)),

		// грейд — только вместе с ролью
		validation.Field(&filter.Grade,
			validation.In("Junior", "Middle", "Senior").Error("грейд: Junior, Middle или Senior"),
			validation.When(filter.Rolle == "", validation.Empty.Error("грейд выбирается только вместе с ролью"))),

		// технологии — не больше 6, без повторов, все есть в таблице stacks
		validation.Field(&filter.StackID,
			validation.Length(0, 6).Error("в поиске не больше 6 технологий"),
			validation.By(StacksExist)),
	)
}
