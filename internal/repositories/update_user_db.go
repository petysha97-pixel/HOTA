package repositories

import (
	"HOTA/internal/models"
	"database/sql"
	"fmt"
)

// обновляем пользователя
func UpdateUser(user models.User, id int) (*models.User, error) {

	qweri := "UPDATE users SET Email = ?, Nickname = ?, Rolle = ?  WHERE id = ?"
	_, err := models.UserDB.Exec(qweri, user.Email, user.Nickname, user.Rolle, id)
	if err != nil {
		return nil, fmt.Errorf("Ошибка в одновлении пользователя %w", err)
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
	qweri := "SELECT id, Email, Nickname, Rolle, about FROM users WHERE id = ?"

	err := models.UserDB.QueryRow(qweri, id).Scan(
		&user.ID,
		&user.Email,
		&user.Nickname,
		&user.Rolle,
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

// Обновляем стеки и обновленные стеки на ручку
func UpdateUserStackID(id int, idstack []int) ([]models.Stack, error) {

	deleteqwery := "DELETE FROM user_stacks WHERE user_id = ?"
	_, err := models.UserDB.Exec(deleteqwery, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка удаления старых секов %w", err)
	}

	if len(idstack) == 0 {
		return nil, fmt.Errorf("Нет стеков для обноваления")
	}

	//сохраненния связей многие ко многим
	for _, stackID := range idstack {
		query := `INSERT INTO user_stacks (user_id, stack_id) VALUES (?, ?)`
		_, err := models.UserDB.Exec(query, id, stackID)

		if err != nil {
			return nil, fmt.Errorf("Ошибка обновления стеков %w", err)
		}

	}

	stacki, err := GetStacksByUserID(id)

	return stacki, nil
}
