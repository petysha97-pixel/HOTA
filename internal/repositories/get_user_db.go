package repositories

import (
	"HOTA/internal/models"
	"database/sql"
	"fmt"
)

// ​​получение пользователя по ИД для отображения данных на главном сайте
func Get_userdb(id int) models.UserResponse {

	var user models.UserResponse
	qwery := "SELECT id, Nickname, COALESCE(Name, ''), Rolle, COALESCE(Grade, ''), COALESCE(about, '') FROM users WHERE id = ?"

	err := models.UserDB.QueryRow(qwery, id).Scan(
		&user.ID,
		&user.Nickname,
		&user.Name,
		&user.Rolle,
		&user.Grade,
		&user.About,
	)
	if err != nil {
		// ИСПРАВЛЕНИЕ: Проверяем, если ошибка — это отсутствие строк
		if err == sql.ErrNoRows {
			// Вместо паники заполняем поля заглушками для фронтенда
			user.Nickname = "Гость"
			user.Rolle = "Не зарегистрирован"
			return user
		}

		// Если это какая-то другая системная ошибка БД — выводим её
		fmt.Println("Критическая ошибка чтения из БД:", err)
		return user
	}

	return user
}

// даостает стеки пользователя для ответа на фронт: название и описание опыта
func GetStacksByUserID(id int) ([]models.Stack, error) {

	stacks := []models.Stack{}

	qwer := `SELECT stacks.id, stacks.name, user_stacks.description FROM user_stacks
	INNER JOIN stacks ON user_stacks.stack_id = stacks.id
	WHERE user_stacks.user_id = ?
	ORDER BY stacks.id`

	rows, err := models.UserDB.Query(qwer, id)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса на получения стека: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stack models.Stack
		err := rows.Scan(&stack.ID, &stack.Name, &stack.Description)
		if err != nil {
			return nil, fmt.Errorf("ошибка записи стка в структуру: %w", err)
		}

		stacks = append(stacks, stack)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения стеков: %w", err)
	}

	return stacks, nil
}
