package project

import (
	"HOTA/internal/models"
	"fmt"
)

func GetProjectPublik() ([]models.Project, error) {
	qweru := `SELECT id, name, description, owner_id, privacy, status, created_at FROM projects 
   WHERE privacy = 'public' 
   AND status IN('draft', 'working')`

	rows, err := models.UserDB.Query(qweru)
	if err != nil {
		return nil, fmt.Errorf("ошибка в получении публичных команат %w", err)

	}
	defer rows.Close()

	var rooms []models.Project
	for rows.Next() {
		var room models.Project
		err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.OwnerID, &room.Privacy, &room.Status, &room.CreatAt)
		if err != nil {
			return nil, fmt.Errorf("Ошибка сканирования комнат %w", err)
		}

		rooms = append(rooms, room)
	}

	return rooms, nil

}
