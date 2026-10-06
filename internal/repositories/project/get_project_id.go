package project

import (
	"HOTA/internal/models"
	"database/sql"
	"fmt"
)

// достаём комнату по айди
func GetProjectByID(roomID int) (*models.Project, error) {
	query := `SELECT id, name, description, owner_id, privacy, status, created_at FROM projects WHERE id = ?`

	var room models.Project
	err := models.UserDB.QueryRow(query, roomID).Scan(&room.ID, &room.Name, &room.Description, &room.OwnerID, &room.Privacy, &room.Status, &room.CreatAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("комната не найдена")
		}
		return nil, fmt.Errorf("ошибка получения комнаты: %w", err)
	}
	return &room, nil
}

// достаём слоты по айди проекта вместе со стеком каждого слота
func GetSlotsByProjectID(projectID int) ([]models.Slot, error) {
	query := `SELECT id, project_id, name, description, rolle, user_id, status, created_at FROM slots WHERE project_id = ?`

	rows, err := models.UserDB.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения слотов: %w", err)
	}
	defer rows.Close()

	slots := []models.Slot{}
	for rows.Next() {
		var slot models.Slot
		var userID sql.NullInt64 // исполнителя может не быть (NULL)

		err := rows.Scan(&slot.ID, &slot.ProjectID, &slot.Name, &slot.Description, &slot.Rolle, &userID, &slot.Status, &slot.CreatAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования слота: %w", err)
		}
		slot.UserID = UserIDOrNil(userID)

		slots = append(slots, slot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения слотов: %w", err)
	}

	// стеки читаем отдельным запросом на каждый слот
	for i := range slots {
		stackIDs, err := GetStackIDsBySlotID(slots[i].ID)
		if err != nil {
			return nil, err
		}
		slots[i].StackID = stackIDs
	}

	return slots, nil
}

// даостает стеки слотов для ответа на фронт
func GetStacksBySlotID(slotID int) ([]models.Stack, error) {

	qwer := `SELECT stacks.id, stacks.name FROM slot_stacks
	INNER JOIN stacks ON slot_stacks.stack_id = stacks.id
	WHERE slot_stacks.slot_id = ?`

	rows, err := models.UserDB.Query(qwer, slotID)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса на получения стека: %w", err)
	}
	defer rows.Close()

	var stacks []models.Stack

	for rows.Next() {
		var stack models.Stack
		if err := rows.Scan(&stack.ID, &stack.Name); err != nil {
			return nil, fmt.Errorf("ошибка записи стека в структуру: %w", err)
		}
		stacks = append(stacks, stack)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения стеков: %w", err)
	}

	return stacks, nil
}

// ID стеков слота — чтобы заполнить Slot.StackID
func GetStackIDsBySlotID(slotID int) ([]int, error) {
	query := `SELECT stack_id FROM slot_stacks WHERE slot_id = ?`

	rows, err := models.UserDB.Query(query, slotID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения стеков слота: %w", err)
	}
	defer rows.Close()

	ids := []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("ошибка сканирования стека слота: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения стеков слота: %w", err)
	}

	return ids, nil
}
