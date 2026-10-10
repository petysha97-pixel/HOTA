package repositories

import (
	"HOTA/internal/models"
	projectRepo "HOTA/internal/repositories/project"
	"fmt"
)

// история участия для профиля:
// 1) проекты, которые пользователь создал
// 2) слоты, где он исполнитель
// tolkoPublichnye = true — для чужого профиля: приватные проекты не показываем
func GetUserHistory(userID int, tolkoPublichnye bool) ([]models.HistoryItem, error) {

	history := []models.HistoryItem{}

	// условие на приватность добавляем только для чужого профиля
	usloviePrivacy := ""
	if tolkoPublichnye {
		usloviePrivacy = " AND projects.privacy = 'public'"
	}

	// 1. созданные проекты + сколько в них слотов всего, занято и сдано
	queryProjects := `SELECT projects.id, projects.name, projects.status, projects.created_at,
		(SELECT COUNT(*) FROM slots WHERE slots.project_id = projects.id),
		(SELECT COUNT(*) FROM slots WHERE slots.project_id = projects.id AND slots.user_id IS NOT NULL AND slots.status != 'done'),
		(SELECT COUNT(*) FROM slots WHERE slots.project_id = projects.id AND slots.status = 'done')
	FROM projects
	WHERE projects.owner_id = ?` + usloviePrivacy + `
	ORDER BY projects.created_at DESC`

	rows, err := models.UserDB.Query(queryProjects, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения проектов пользователя: %w", err)
	}

	for rows.Next() {
		item := models.HistoryItem{Type: "project"}
		err := rows.Scan(&item.ProjectID, &item.ProjectName, &item.Status, &item.CreatAt, &item.SlotsTotal, &item.SlotsTaken, &item.SlotsDone)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("ошибка сканирования проекта: %w", err)
		}
		history = append(history, item)
	}
	rows.Close()

	// 2. слоты, где пользователь исполнитель + проект и ник его создателя
	querySlots := `SELECT slots.id, slots.name, slots.rolle, slots.status, slots.created_at,
		projects.id, projects.name, users.Nickname
	FROM slots
	INNER JOIN projects ON slots.project_id = projects.id
	INNER JOIN users ON projects.owner_id = users.id
	WHERE slots.user_id = ?` + usloviePrivacy + `
	ORDER BY slots.created_at DESC`

	rows, err = models.UserDB.Query(querySlots, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения слотов пользователя: %w", err)
	}

	slotsStart := len(history) // с этого места в списке начинаются слоты
	for rows.Next() {
		item := models.HistoryItem{Type: "slot"}
		err := rows.Scan(&item.SlotID, &item.SlotName, &item.Rolle, &item.Status, &item.CreatAt,
			&item.ProjectID, &item.ProjectName, &item.OwnerNickname)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("ошибка сканирования слота: %w", err)
		}
		history = append(history, item)
	}
	rows.Close()

	// стек каждого слота читаем отдельно (rows уже закрыт, иначе SQLite занят)
	for i := slotsStart; i < len(history); i++ {
		stackIDs, err := projectRepo.GetStackIDsBySlotID(history[i].SlotID)
		if err != nil {
			return nil, err
		}
		history[i].StackID = stackIDs
	}

	return history, nil
}
