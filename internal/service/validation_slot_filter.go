package service

import (
	"HOTA/internal/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// проверка фильтров вкладки «Слоты»
func ValidateSlotFilter(filter models.SlotFilter) error {
	return validation.ValidateStruct(&filter,

		// вкладка
		validation.Field(&filter.Tab, validation.In("open", "mine").Error("вкладка: open или mine")),

		// текст поиска
		validation.Field(&filter.Query, validation.Length(0, 50).Error("текст поиска не должен быть больше 50 символов")),

		// роль — из каталога ролей
		validation.Field(&filter.Rolle, validation.By(RoleExists)),

		// технологии — не больше 6, без повторов, все есть в таблице stacks
		validation.Field(&filter.StackID,
			validation.Length(0, 6).Error("в поиске не больше 6 технологий"),
			validation.By(StacksExist)),
	)
}
