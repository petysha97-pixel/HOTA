package project

import (
	"HOTA/internal/models"
	"database/sql"
	"fmt"
)

// левая панель «Проекты»: что я создал и где я в команде
func GetMyProjects(meID int) (models.MyProjectsResponse, error) {

	var response models.MyProjectsResponse
	var err error

	// «Создал»
	response.Created, err = getRailProjects(`SELECT id, name, status FROM projects
		WHERE owner_id = ? ORDER BY id DESC`, meID)
	if err != nil {
		return response, err
	}

	// «Участвую»: я исполнитель хотя бы в одном слоте чужого проекта
	response.Participate, err = getRailProjects(`SELECT DISTINCT projects.id, projects.name, projects.status
		FROM projects
		INNER JOIN slots ON slots.project_id = projects.id
		WHERE slots.user_id = ? AND projects.owner_id != ?
		ORDER BY projects.id DESC`, meID, meID)
	if err != nil {
		return response, err
	}

	// слоты, мои слоты и цифры бейджа
	for i := range response.Created {
		if err := fillRailProject(&response.Created[i], meID, true); err != nil {
			return response, err
		}
	}
	for i := range response.Participate {
		if err := fillRailProject(&response.Participate[i], meID, false); err != nil {
			return response, err
		}
	}

	return response, nil
}

// список проектов панели (id, название, статус) по готовому запросу
func getRailProjects(query string, args ...any) ([]models.RailProject, error) {

	rows, err := models.UserDB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения моих проектов: %w", err)
	}
	defer rows.Close()

	projects := []models.RailProject{}
	for rows.Next() {
		var p models.RailProject
		if err := rows.Scan(&p.ID, &p.Name, &p.Status); err != nil {
			return nil, fmt.Errorf("ошибка сканирования проекта: %w", err)
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения проектов: %w", err)
	}

	return projects, nil
}

// дозаполняем проект панели: слоты, мои слоты, цифры бейджа
// isOwner — я создатель проекта (тогда считаем отклики и review)
func fillRailProject(p *models.RailProject, meID int, isOwner bool) error {

	rows, err := models.UserDB.Query(`SELECT id, name, status, user_id FROM slots WHERE project_id = ? ORDER BY id`, p.ID)
	if err != nil {
		return fmt.Errorf("ошибка получения слотов: %w", err)
	}
	defer rows.Close()

	p.Slots = []models.RailSlot{}
	p.MySlots = []string{}
	for rows.Next() {
		var slot models.RailSlot
		var userID sql.NullInt64 // исполнителя может не быть (NULL)

		if err := rows.Scan(&slot.ID, &slot.Name, &slot.Status, &userID); err != nil {
			return fmt.Errorf("ошибка сканирования слота: %w", err)
		}
		slot.UserID = UserIDOrNil(userID)
		p.Slots = append(p.Slots, slot)

		// занят: есть исполнитель, работа ещё не сдана
		if slot.UserID != nil && slot.Status != "done" {
			p.SlotsTaken++
		}
		// мой слот
		if slot.UserID != nil && *slot.UserID == meID {
			p.MySlots = append(p.MySlots, slot.Name)
		}
		// на знакомстве — показываем только создателю
		if isOwner && slot.Status == "review" {
			p.Review++
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("ошибка чтения слотов: %w", err)
	}
	p.SlotsTotal = len(p.Slots)

	// отклики на слоты проекта, которые ждут ответа — только создателю
	if isOwner {
		err = models.UserDB.QueryRow(`SELECT COUNT(*) FROM applicationsSlot
			INNER JOIN slots ON slots.id = applicationsSlot.slot_id
			WHERE slots.project_id = ? AND applicationsSlot.status = 'pending'`, p.ID).Scan(&p.Applications)
		if err != nil {
			return fmt.Errorf("ошибка подсчёта откликов: %w", err)
		}
	}

	// приглашения мне в этот проект, которые ждут ответа
	err = models.UserDB.QueryRow(`SELECT COUNT(*) FROM invites
		INNER JOIN slots ON slots.id = invites.slot_id
		WHERE slots.project_id = ? AND invites.user_id = ? AND invites.status = 'pending'`, p.ID, meID).Scan(&p.Invites)
	if err != nil {
		return fmt.Errorf("ошибка подсчёта приглашений: %w", err)
	}

	p.Badge = p.Applications + p.Review + p.Invites
	return nil
}
