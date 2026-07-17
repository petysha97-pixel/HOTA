package repositories

import (
	"HOTA/internal/models"
	"strings"
)


func SearchUser(searchuser string) ([]models.User, error) {

	sqlQuery := `
		SELECT id, Nickname, Role, Stack
		FROM users
		WHERE Nickname LIKE ?
		   OR Role LIKE ?
		   OR Stack LIKE ?
	`

	rows, err := models.UserDB.Query(
		sqlQuery,

		// Ищем совпадение в нике
		searchuser+"%",

		// Ищем совпадение в роли
		"%"+searchuser+"%",

		// Ищем совпадение в стеке
		"%"+searchuser+"%",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	
	var users []models.User

	
	for rows.Next() {

		var user models.User
		var stackStr string

		
		err := rows.Scan(
			&user.Nickname,
			&user.Rolle,
			&stackStr,
		)

		if err != nil {
			return nil, err
		}

		
		if stackStr != "" {
			user.Stack = strings.Split(stackStr, ",")
		}

		
		users = append(users, user)
	}

	
	return users, err
}