package repositories

import (
	"HOTA/internal/models"
	"fmt"
)

// отдаём весь каталог ролей (для регистрации, профиля и слотов)
func GetRoles() ([]models.Role, error) {

	rows, err := models.UserDB.Query(`SELECT id, name FROM roles ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ролей: %w", err)
	}
	defer rows.Close()

	roles := []models.Role{}
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, fmt.Errorf("ошибка сканирования роли: %w", err)
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения ролей: %w", err)
	}

	return roles, nil
}
