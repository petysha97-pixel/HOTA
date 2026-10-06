package service

import (
	"HOTA/internal/models"
	"database/sql"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// валидация полей при регистрации
func ValidateUpdataStruct(User models.User, UserID int) error {
	return validation.ValidateStruct(&User,

		//Почта (мыло @ + домен)
		validation.Field(&User.Email, validation.Required, is.Email, validation.Length(1, 30), validation.By(UNIL_email_updata(UserID))),

		// Никнейм (от 2 символов)
		validation.Field(&User.Nickname, validation.Required, validation.Length(2, 20), validation.By(UNIL_nikaname_updata(UserID))),

		// Роль обязательна и должна быть одной из строго заданных на фронтенде
		validation.Field(&User.Rolle, validation.Required, validation.By(RoleExists)),

		// Грейд обязателен
		validation.Field(&User.Grade, validation.Required, validation.In(Grades...)),

		// ФИО по желанию
		validation.Field(&User.Name, validation.Length(0, 60)),

		// Стек обязателен. Мы проверяем каждый элемент массива (каждую строку технологии)
		validation.Field(&User.StackID, validation.Required, validation.Length(1, 6), validation.Each(
			validation.Required, // Минимум 1 стек
		), validation.By(StacksExist)),
	)
}

func UNIL_email_updata(UserID int) func(any) error {
	return func(email any) error {

		emal, ok := email.(string)
		if !ok {
			return errors.New("Почта должна быть строкой")
		}

		var count int

		query := `SELECT COUNT(1) FROM users WHERE Email = ? AND id != ?`

		err := models.UserDB.QueryRow(query, emal, UserID).Scan(&count)

		if errors.Is(err, sql.ErrNoRows) {
			// логин свободен
			return nil
		}

		if err != nil {
			return err //ошиба из БД
		}

		if count > 0 {
			return errors.New("емаил уже занят другим пользователем")
		}

		return nil
	}

}

func UNIL_nikaname_updata(UserID int) func(any) error {
	return func(nikaname any) error {

		nik, ok := nikaname.(string)
		if !ok {
			return errors.New("Ник должен быть строкой")
		}

		var count int

		query := `SELECT COUNT(1) FROM users WHERE Nickname = ? AND id != ?`

		err := models.UserDB.QueryRow(query, nik, UserID).Scan(&count)

		if errors.Is(err, sql.ErrNoRows) {
			// ник свободен
			return nil
		}

		if err != nil {
			return err //ошиба из БД
		}

		if count > 0 {
			return errors.New("ник уже занят другим пользователем")
		}

		return nil
	}

}
