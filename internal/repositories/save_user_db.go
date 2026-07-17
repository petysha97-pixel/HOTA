package repositories

import (
	"HOTA/internal/models"
	"database/sql"
)

// AppendUser добавляет нового пользователя в БД и возвращает его же с ID
func AppendUser(user models.User) (models.User, error) {

	for _, tech := range user.Stack {

		var count int

		query := `SELECT 1 FROM stacks WHERE name = ? LIMIT 1`

		err := models.UserDB.QueryRow(query, tech).
			Scan(&count)

		if err == sql.ErrNoRows {
			return models.User{}, nil
		}

		if err != nil {
			return models.User{}, err
		}
	}

	userSQL := `
	INSERT INTO users (
		Login,
		Password,
	    Nickname,
		Rolle
	   )
	VALUES (?, ?, ?, ?);
	`

	row, err := models.UserDB.Exec(
		userSQL,
		user.Login,
		user.Password,
		user.Nickname,
		user.Rolle,
	)
	if err != nil {
		return models.User{}, err
	}

	//достаем id пользователя из БД
	id, _ := row.LastInsertId()
	user.ID = int(id)

	for _, tech := range user.Stack {
		var stackID int
		query := `SELECT id FROM stacks WHERE name = ?`
		err := models.UserDB.QueryRow(query, tech).Scan(&stackID)
		if err != nil {
			return models.User{}, err
		}

		query = `INSERT INTO user_stacks (user_id, stack_id) SELECT 2, 2 WHERE EXISTS (SELECT 1 FROM stacks WHERE id = 2);`
		_, err = models.UserDB.Exec(query, User.ID, stackID, stackID)

		if err != nil {
			return models.User{}, err
		}

	}

	return user, err
}
