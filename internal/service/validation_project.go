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

		//Описание проекта
		validation.Field(&project.Description, validation.Required.Error("Описание проекта не может быть пустым"),
			validation.Length(15, 300).Error("Описание проекта не должно быть меньше 15 символов и больше 300")),

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

		//название слота
		validation.Field(&slot.Name, validation.Required.Error("Название слота не может быть пустым"),
			validation.Length(3, 50).Error("Название слота не должно быть меньше 3 символов и больше 50")),

		//описание слота (техзадание) — по желанию
		validation.Field(&slot.Description, validation.Length(0, 300).Error("Описание слота не должно быть больше 300 символов")),

		//роль слота — из каталога ролей
		validation.Field(&slot.Rolle, validation.Required, validation.By(RoleExists)),

		//стеки слота
		validation.Field(&slot.StackID, validation.Required, validation.Length(1, 6), validation.Each(
			validation.Required, // Минимум 1 стек
		), validation.By(StacksExist)),
	)
}
