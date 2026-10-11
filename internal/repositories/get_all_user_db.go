package repositories

import (
	"HOTA/internal/models"
	"fmt"
	"strings"
)

// разработчики для вкладки «Разработчики»
// meID — кто ищет: себя в списке не показываем
// фильтры роль, грейд и стек делаем в SQL, текст ищем в хендлере
func GetUsers(filter models.UserFilter, meID int) ([]models.User, error) {

	query := `SELECT id, Nickname, COALESCE(Name, ''), Rolle, COALESCE(Grade, ''), COALESCE(about, '')
	FROM users WHERE id != ?`
	args := []any{meID}

	// роль
	if filter.Rolle != "" {
		query += ` AND Rolle = ?`
		args = append(args, filter.Rolle)
	}

	// грейд
	if filter.Grade != "" {
		query += ` AND Grade = ?`
		args = append(args, filter.Grade)
	}

	// стек: у разработчика должны быть ВСЕ выбранные технологии
	// считаем, сколько из выбранных у него есть, и сравниваем с количеством выбранных
	if len(filter.StackID) > 0 {
		voprosy := strings.TrimSuffix(strings.Repeat("?,", len(filter.StackID)), ",") // "?,?,?"
		query += ` AND id IN (
			SELECT user_id FROM user_stacks
			WHERE stack_id IN (` + voprosy + `)
			GROUP BY user_id
			HAVING COUNT(*) = ?)`
		for _, id := range filter.StackID {
			args = append(args, id)
		}
		args = append(args, len(filter.StackID))
	}

	// новые разработчики сверху
	query += ` ORDER BY id DESC`

	rows, err := models.UserDB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса разработчиков: %w", err)
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Nickname, &user.Name, &user.Rolle, &user.Grade, &user.About)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования разработчика: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения разработчиков: %w", err)
	}

	return users, nil
}