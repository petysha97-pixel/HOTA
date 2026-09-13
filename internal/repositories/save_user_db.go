package repositories

import (
	"HOTA/internal/models"
)

// AppendUser добавляет нового пользователя в БД и возвращает его же с ID
func AppendUser(user models.User) (models.User, error) {

	// for _, tech := range user.Stack {

	// 	var count int

	// 	query := `SELECT 1 FROM stacks WHERE name = ? LIMIT 1`

	// 	err := models.UserDB.QueryRow(query, tech).
	// 		Scan(&count)

	// 	if err == sql.ErrNoRows {
	// 		return models.User{}, nil
	// 	}

	// 	if err != nil {
	// 		return models.User{}, err
	// 	}
	// }

	userSQL := `
	INSERT INTO users (
		Email,
		Password,
	    Nickname,
		Rolle
	   )
	VALUES (?, ?, ?, ?);
	`

	row, err := models.UserDB.Exec(
		userSQL,
		user.Email,
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

	//сохраненния связей многие ко многим
	for _, stackID := range user.Stack {
		query := `INSERT INTO user_stacks (user_id, stack_id) VALUES (?, ?)`
		_, err := models.UserDB.Exec(query, user.ID, stackID)

		if err != nil {
			return models.User{}, err
		}

	}

	return user, err
}
