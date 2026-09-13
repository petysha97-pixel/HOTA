package service

import (
	"HOTA/internal/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// валидация проекта при регистрации
func ValidateProlectStruct(project models.DTOProject) error {
	return validation.ValidateStruct(&project,

		//Название проекта
		validation.Field(&project.Name, validation.Required.Error("Название проекта не может быть пустым"),
			validation.Length(3, 50).Error("Название проекта не должно быть меньше 3 символов и больше 50")),

		//Цель проекта
		validation.Field(&project.Target, validation.Required.Error("Цель проекта не может быть пустым"),
			validation.Length(15, 300).Error("Цель проекта не должна быть меньше 15 символов и больше 300")),

		//приватность проекта
		validation.Field(&project.Privacy, validation.In("public", "private").Error("Укажите приватность проекта")),

		//статус проекта
		validation.Field(&project.Status, validation.In("draft", "working", "finished").Error("Укажите статус проекта")),

		//слоты проекта
		validation.Field(&project.Slots, validation.Required.Error("Должен быть хотя бы 1 слот в проекте"),
			validation.Length(1, 6).Error("Минимум 1 слот, максимум 6")),
	)

}

// валидация слота при регистрации проета
func ValidateSlotStruct(slot models.DTOSlot) error {
	return validation.ValidateStruct(&slot,

		//роль слота
		validation.Field(&slot.Rolle, validation.Required, validation.In(
			"Frontend", "Backend", "Fullstack", "DevOps")),

		//статус слота
		validation.Field(&slot.Status, validation.In("open", "close", "done").Error("Укажите статус слота")),

		//стеки слота
		validation.Field(&slot.Stack, validation.Required, validation.Length(1, 6), validation.Each(
			validation.Required, // Минимум 1 стек
		)),
	)
}
