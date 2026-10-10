package service

import (
	"HOTA/internal/models"
	"database/sql"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// валидация правки профиля (окно «Профиль»): ник, ФИО, роль, грейд, о себе
// почта и стек меняются своими ручками и проверяются там
func ValidateUpdataStruct(profile models.UpdateProfile, UserID int) error {
	return validation.ValidateStruct(&profile,

		// Никнейм (от 2 символов), свободный у других пользователей
		validation.Field(&profile.Nickname, validation.Required, validation.Length(2, 20), validation.By(UNIL_nikaname_updata(UserID))),

		// ФИО по желанию
		validation.Field(&profile.Name, validation.Length(0, 60)),

		// Роль обязательна и должна быть в каталоге ролей
		validation.Field(&profile.Rolle, validation.Required, validation.By(RoleExists)),

		// Грейд обязателен
		validation.Field(&profile.Grade, validation.Required, validation.In(Grades...)),

		// О себе по желанию, до 200 символов
		validation.Field(&profile.About, validation.Length(0, 200)),
	)
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
