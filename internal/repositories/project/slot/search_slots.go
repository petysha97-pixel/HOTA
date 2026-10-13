package slot

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	"database/sql"
	"fmt"
	"strings"
)

// слоты для вкладки «Слоты»
// берём только слоты, куда ещё можно попасть:
//   - проект публичный, не завершён и не мой
//   - слот open (ищут человека) или review (идёт знакомство)
//
// фильтры роль, «только свободные» и стек делаем в SQL, остальное — в хендлере
func SearchSlots(filter models.SlotFilter, meID int) ([]models.SlotCard, error) {

	query := `SELECT slots.id, slots.project_id, projects.name, slots.name, slots.description,
	slots.rolle, slots.status, slots.user_id
	FROM slots
	INNER JOIN projects ON projects.id = slots.project_id
	WHERE projects.privacy = 'public'
	AND projects.status != 'finished'
	AND projects.owner_id != ?
	AND slots.status IN ('open', 'review')`
	args := []any{meID}

	// роль слота
	if filter.Rolle != "" {
		query += ` AND slots.rolle = ?`
		args = append(args, filter.Rolle)
	}

	// только свободные: без кандидата на знакомстве
	if filter.OnlyFree {
		query += ` AND slots.status = 'open'`
	}

	// стек: слоту нужны ВСЕ выбранные технологии
	if len(filter.StackID) > 0 {
		voprosy := strings.TrimSuffix(strings.Repeat("?,", len(filter.StackID)), ",") // "?,?,?"
		query += ` AND slots.id IN (
			SELECT slot_id FROM slot_stacks
			WHERE stack_id IN (` + voprosy + `)
			GROUP BY slot_id
			HAVING COUNT(*) = ?)`
		for _, id := range filter.StackID {
			args = append(args, id)
		}
		args = append(args, len(filter.StackID))
	}

	// новые слоты сверху
	query += ` ORDER BY slots.id DESC`

	rows, err := models.UserDB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса слотов: %w", err)
	}
	defer rows.Close()

	slots := []models.SlotCard{}
	for rows.Next() {
		var slot models.SlotCard
		var userID sql.NullInt64 // кандидата может не быть (NULL)

		err := rows.Scan(&slot.ID, &slot.ProjectID, &slot.ProjectName, &slot.Name, &slot.Description,
			&slot.Rolle, &slot.Status, &userID)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования слота: %w", err)
		}
		slot.UserID = projectRepo.UserIDOrNil(userID)

		slots = append(slots, slot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения слотов: %w", err)
	}

	// стек каждого слота (id и название)
	for i := range slots {
		stack, err := projectRepo.GetStacksBySlotID(slots[i].ID)
		if err != nil {
			return nil, err
		}
		if stack == nil {
			stack = []models.Stack{}
		}
		slots[i].Stack = stack
	}

	return slots, nil
}

// id слотов, где мой отклик ждёт ответа (status = pending)
func GetMyPendingSlotIDs(meID int) (map[int]bool, error) {

	rows, err := models.UserDB.Query(`SELECT slot_id FROM applicationsSlot WHERE user_id = ? AND status = 'pending'`, meID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения моих откликов: %w", err)
	}
	defer rows.Close()

	slotIDs := map[int]bool{}
	for rows.Next() {
		var slotID int
		if err := rows.Scan(&slotID); err != nil {
			return nil, fmt.Errorf("ошибка сканирования отклика: %w", err)
		}
		slotIDs[slotID] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения откликов: %w", err)
	}

	return slotIDs, nil
}
