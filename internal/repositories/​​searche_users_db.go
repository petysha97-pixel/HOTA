package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

// ищет пользавателей по запросу юзера с главного меню
func SearcheUsersBD(query, limit string) ([]models.User, error) {

	rules := "%" + query + "%" //убрать первый % для того, что поиск не выдавал все слова по типу "Django"
	//нужно написать логику поиска разрабочикjd по многим составляющим (go+backend+petya)
	// + толерантность к ошибкам (через температуру схожести) + приводиться слова к 1 региситру

	quer := `SELECT DISTINCT users.id, users.Nickname, users.Rolle FROM users 
	LEFT JOIN user_stacks
	ON users.id = user_stacks.user_id
	LEFT JOIN stacks 
	ON user_stacks.stack_id = stacks.id
	WHERE Nickname LIKE ?
	OR users.Rolle LIKE ?
	OR stacks.name LIKE ?
	  ORDER BY
	       CASE 
		     WHEN users.Nickname LIKE ? THEN 1 
			 WHEN stacks.name LIKE ? THEN 2
			 ELSE 3
			 END
			 LIMIT ?`

	rows, err := models.UserDB.Query(quer, rules, rules, rules, rules, rules, limit)
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса в БД: SearcheUsersBD %w", err)
	}

	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Nickname, &user.Rolle)
		if err != nil {
			return nil, fmt.Errorf("Ошибка сканирования пользователя из БД %w", err)
		}

		users = append(users, user)
	}

	return users, nil

}
