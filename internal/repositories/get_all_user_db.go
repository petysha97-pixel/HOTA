package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

// функция короторая берет всех пользователей из БД во фронт
func GETUser() ([]models.User, error) {

	квери := "SELECT id, Nickname, COALESCE(Name, ''), Rolle, COALESCE(Grade, '') FROM users"

	ровс, ошибка := models.UserDB.Query(квери)
	if ошибка != nil {
		return nil, fmt.Errorf("Ошибка запроса в БД: SearchUser %w", ошибка)
	}
	//МНЕ ПОДСКАЗАЛ LAPA (ТАК ВСЕ ДЕЛАЮТ)

	defer ровс.Close()

	var users []models.User
	for ровс.Next() {
		var user models.User
		err := ровс.Scan(&user.ID, &user.Nickname, &user.Name, &user.Rolle, &user.Grade)
		if err != nil {
			return nil, fmt.Errorf("Ошибка сканирования пользователя из БД %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}
