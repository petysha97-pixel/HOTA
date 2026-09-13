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

// достаём слоты по айди комнаты
func GetSlotsByProjectID(projectID int) ([]models.Slot, error) {
	query := `SELECT id, project_id, rolle, status, created_at FROM slots WHERE project_id = ?`

	rows, err := models.UserDB.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения слотов: %w", err)
	}
	defer rows.Close()

	var slots []models.Slot
	for rows.Next() {
		var slot models.Slot
		err := rows.Scan(&slot.ID, &slot.ProjectID, &slot.Rolle, &slot.Status, &slot.CreatAt)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования слота: %w", err)
		}
		slots = append(slots, slot)
	}


	
	fmt.Println(slots)
	return slots, nil
}



//даостает стеки слотов для ответа на фронт
func GetStacksBySlotID(id int) ([]models.Stack, error){
    
	
	var stacks []models.Stack
	
	qwer := `SELECT stacks.id, stacks.name FROM user_stacks 
	INNER JOIN stacks ON user_stacks.stack_id = stacks.id
	WHERE user_stacks.user_id = ?`

	rows, err := models.UserDB.Query(qwer, id)
    if err != nil{
		return nil, fmt.Errorf("ошибка выполнения запроса на получения стека: %w", err)
	}
	defer rows.Close()
	
	
	for rows.Next() {
  	var stack models.Stack
	err := rows.Scan(&stack.ID, &stack.Name)
    if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка записи стка в структуру: %w", err)
	}

	stacks = append(stacks, stack)
	}


	return stacks, nil
}
