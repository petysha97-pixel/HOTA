package project

import (
	"HOTA/internal/models"
	"fmt"
)

// создние комнаты
func CreatProject(room *models.Project) error {
	query := `INSERT INTO projects (name, description, owner_id, privacy, status)VALUES (?, ?, ?, ?, ?);`

	row, err := models.UserDB.Exec(query, room.Name, room.Description, room.OwnerID, room.Privacy, room.Status)
	if err != nil {
		return fmt.Errorf("ошибка создание комнаты %w", err)
	}

	//достаем id комнаты из БД
	id, _ := row.LastInsertId()
	room.ID = int(id)

	return nil
}
