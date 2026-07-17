package repositories

import (
	"HOTA/internal/models"
	"database/sql"
	"fmt"
)

// ​​получение пользователя по ИД для отображения данных на главном сайте
func Get_userdb(id int) models.User {

	var user models.User
	qwery := "SELECT id, Nickname, Rolle FROM users WHERE id = ?"

	err := models.UserDB.QueryRow(qwery, id).Scan(
		&user.ID,
		&user.Nickname,
		&user.Rolle,
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
