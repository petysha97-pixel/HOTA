package project

import (
	"HOTA/internal/models"
	"database/sql"
	"fmt"
)

// записываем слот и его стек
// вызываем только внутри транзакции (tx), чтобы при ошибке всё откатилось
func createSlot(tx *sql.Tx, slot *models.Slot) error {

	// id и дату создания ставит сама БД, забираем их через RETURNING
	query := `INSERT INTO slots (project_id, name, description, rolle, status) VALUES (?, ?, ?, ?, ?) RETURNING id, created_at`

	err := tx.QueryRow(query, slot.ProjectID, slot.Name, slot.Description, slot.Rolle, slot.Status).Scan(&slot.ID, &slot.CreatAt)
	if err != nil {
		return fmt.Errorf("ошибка создание слота %w", err)
	}

	//сохраненния связей многие ко многим (слоту к стекам)
	for _, stackID := range slot.StackID {
		_, err := tx.Exec(`INSERT INTO slot_stacks (slot_id, stack_id) VALUES (?, ?)`, slot.ID, stackID)
		if err != nil {
			return fmt.Errorf("ошибка добавления стека в слот %w", err)
		}
	}

	return nil
}

// создание слота вместе со стеком (ручка POST /project/{id}/slot)
func CreateSlot(slot *models.Slot) error {

	// начинаем транзакцию
	tx, err := models.UserDB.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback() // если не дошли до Commit — всё отменится

	err = createSlot(tx, slot)
	if err != nil {
		return err
	}

	// сохраняем изменения
	return tx.Commit()
}
