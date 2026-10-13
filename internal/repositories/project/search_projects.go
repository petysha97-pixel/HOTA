package project

import (
	"HOTA/internal/models"
	"database/sql"
	"fmt"
)

// проекты для вкладки «Проекты»
// видимость: публичный проект, или я создатель, или я исполнитель в его слоте
// фильтры статус, «Мои» и «Моя роль» делаем в SQL, текст ищем в хендлере
func SearchProjects(filter models.ProjectFilter, meID int) ([]models.ProjectCard, error) {

	query := `SELECT projects.id, projects.name, projects.description, projects.status, projects.privacy,
	projects.created_at, projects.owner_id, users.Nickname
	FROM projects
	INNER JOIN users ON users.id = projects.owner_id
	WHERE (projects.privacy = 'public'
		OR projects.owner_id = ?
		OR EXISTS (SELECT 1 FROM slots WHERE slots.project_id = projects.id AND slots.user_id = ?))`
	args := []any{meID, meID}

	// статус: завершённые видны только при status=finished
	if filter.Status != "" {
		query += ` AND projects.status = ?`
		args = append(args, filter.Status)
	} else {
		query += ` AND projects.status != 'finished'`
	}

	// working показываем, только если в нём ещё можно занять слот (open или review)
	query += ` AND (projects.status != 'working'
		OR EXISTS (SELECT 1 FROM slots WHERE slots.project_id = projects.id AND slots.status IN ('open', 'review')))`

	// «Мои»: создал сам или исполнитель в слоте
	if filter.Filter == "mine" {
		query += ` AND (projects.owner_id = ?
			OR EXISTS (SELECT 1 FROM slots WHERE slots.project_id = projects.id AND slots.user_id = ?))`
		args = append(args, meID, meID)
	}

	// «Моя роль»: есть open или review слот моей роли
	if filter.Filter == "role" {
		query += ` AND EXISTS (SELECT 1 FROM slots WHERE slots.project_id = projects.id
			AND slots.rolle = ? AND slots.status IN ('open', 'review'))`
		args = append(args, filter.Rolle)
	}

	// новые проекты сверху
	query += ` ORDER BY projects.id DESC`

	rows, err := models.UserDB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса проектов: %w", err)
	}
	defer rows.Close()

	projects := []models.ProjectCard{}
	for rows.Next() {
		var p models.ProjectCard
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.Privacy, &p.CreatAt, &p.OwnerID, &p.OwnerNickname)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования проекта: %w", err)
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения проектов: %w", err)
	}

	// слоты каждого проекта и счётчики по статусам
	for i := range projects {
		slots, err := GetSlotsShortByProjectID(projects[i].ID)
		if err != nil {
			return nil, err
		}
		projects[i].Slots = slots
		projects[i].Counts = CountSlots(slots)
	}

	return projects, nil
}

// один проект для страницы проекта (вместе со слотами и ником автора)
func GetProjectCard(projectID int) (*models.ProjectCard, error) {

	query := `SELECT projects.id, projects.name, projects.description, projects.status, projects.privacy,
	projects.created_at, projects.owner_id, users.Nickname
	FROM projects
	INNER JOIN users ON users.id = projects.owner_id
	WHERE projects.id = ?`

	var p models.ProjectCard
	err := models.UserDB.QueryRow(query, projectID).Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.Privacy, &p.CreatAt, &p.OwnerID, &p.OwnerNickname)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("проект не найден")
		}
		return nil, fmt.Errorf("ошибка получения проекта: %w", err)
	}

	p.Slots, err = GetSlotsShortByProjectID(p.ID)
	if err != nil {
		return nil, err
	}
	p.Counts = CountSlots(p.Slots)

	return &p, nil
}

// слоты проекта со стеком (id и название технологий)
func GetSlotsShortByProjectID(projectID int) ([]models.SlotShort, error) {

	query := `SELECT id, name, description, rolle, status, user_id FROM slots WHERE project_id = ? ORDER BY id`

	rows, err := models.UserDB.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения слотов: %w", err)
	}
	defer rows.Close()

	slots := []models.SlotShort{}
	for rows.Next() {
		var slot models.SlotShort
		var userID sql.NullInt64 // исполнителя может не быть (NULL)

		err := rows.Scan(&slot.ID, &slot.Name, &slot.Description, &slot.Rolle, &slot.Status, &userID)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования слота: %w", err)
		}
		slot.UserID = UserIDOrNil(userID)

		slots = append(slots, slot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения слотов: %w", err)
	}

	// стек каждого слота
	for i := range slots {
		stack, err := GetStacksBySlotID(slots[i].ID)
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

// сколько слотов в каждом статусе
func CountSlots(slots []models.SlotShort) models.ProjectCounts {
	var c models.ProjectCounts
	for _, slot := range slots {
		switch slot.Status {
		case "open":
			c.Open++
		case "review":
			c.Review++
		case "close":
			c.Close++
		case "done":
			c.Done++
		}
	}
	c.Total = len(slots)
	return c
}
