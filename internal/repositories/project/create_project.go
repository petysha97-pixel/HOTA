package project

import (
	"HOTA/internal/models"
	"fmt"
)

// создание проекта вместе со слотами одной транзакцией: либо всё, либо ничего
func CreateProjectWithSlots(room *models.Project, slots []models.Slot) error {
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() // откатится, если не закоммитим

	// id и дату создания ставит БД, забираем их через RETURNING
	query := `INSERT INTO projects (name, description, owner_id, privacy, status) VALUES (?, ?, ?, ?, ?) RETURNING id, created_at`

	err = tx.QueryRow(query, room.Name, room.Description, room.OwnerID, room.Privacy, room.Status).Scan(&room.ID, &room.CreatAt)
	if err != nil {
		return fmt.Errorf("ошибка создание комнаты %w", err)
	}

	for i := range slots {
		slots[i].ProjectID = room.ID
		if err := createSlot(tx, &slots[i]); err != nil {
			return err
		}
	}

	return tx.Commit()
}
