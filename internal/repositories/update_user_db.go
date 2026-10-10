package repositories

import (
	"HOTA/internal/models"
	"database/sql"
	"fmt"
)

// обновляем профиль: ник, ФИО, роль, грейд, о себе
// почта и стек меняются своими ручками
func UpdateUser(profile models.UpdateProfile, id int) error {

	qweri := "UPDATE users SET Nickname = ?, Name = ?, Rolle = ?, Grade = ?, about = ?, Update_add = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := models.UserDB.Exec(qweri, profile.Nickname, profile.Name, profile.Rolle, profile.Grade, profile.About, id)
	if err != nil {
		return fmt.Errorf("Ошибка в одновлении пользователя %w", err)
	}

	return nil
}

// достает пользователя из БД по айди
func GetUsersByID(id int) (*models.User, error) {

	var user models.User
	var about sql.NullString
	qweri := "SELECT id, Email, Nickname, COALESCE(Name, ''), Rolle, COALESCE(Grade, ''), about FROM users WHERE id = ?"

	err := models.UserDB.QueryRow(qweri, id).Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.Name,
		&user.Rolle,
		&user.Grade,
		&about,
	)

	if err != nil {
		return nil, fmt.Errorf("Пользователя с айди %d не существует: %w", id, err)
	}

	if about.Valid {
		user.About = about.String
	} else {
		user.About = ""
	}

	return &user, nil

}
