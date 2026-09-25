package project

import (
	"HOTA/internal/models"
	"fmt"
)

// создание слота
func CreateSlot(slot *models.Slot) error {

	query := `INSERT INTO slots (project_id, rolle, status) VALUES (?, ?, ?);`

	row, err := models.UserDB.Exec(query, slot.ProjectID, slot.Rolle, slot.Status)
	if err != nil {
		return fmt.Errorf("ошибка создание слота %w", err)
	}

	//достаем id слота из БД
	id, _ := row.LastInsertId()
	slot.ID = int(id)

	//сохраненния связей многие ко многим (слоту к стекам)
	for _, stackID := range slot.StackID {
		query := `INSERT INTO slot_stacks (slot_id, stack_id) VALUES (?, ?)`
		_, err := models.UserDB.Exec(query, slot.ID, stackID)

		if err != nil {
			return fmt.Errorf("ошибка добавления стека в слот %w", err)
		}

	}

	

	return nil
}
