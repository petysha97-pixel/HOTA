package repositories

import (
	"HOTA/internal/models"
	"database/sql"
	"fmt"
)

// обновляем пользователя вместе со стеком одной транзакцией
func UpdateUser(user models.User, id int) (*models.User, error) {

	tx, err := models.UserDB.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() // откатится, если не закоммитим

	qweri := "UPDATE users SET Email = ?, Nickname = ?, Name = ?, Rolle = ?, Grade = ?, Update_add = CURRENT_TIMESTAMP WHERE id = ?"
	_, err = tx.Exec(qweri, user.Email, user.Nickname, user.Name, user.Rolle, user.Grade, id)
	if err != nil {
		return nil, fmt.Errorf("Ошибка в одновлении пользователя %w", err)
	}

	if err := replaceUserStacks(tx, id, user.StackID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	newuser, err := GetUsersByID(id)
	if err != nil {
		return nil, fmt.Errorf("Ошибка в одновлении пользователя %w", err)
	}

	return newuser, nil
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
