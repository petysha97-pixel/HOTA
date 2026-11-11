package slot

import (
	"HOTA/internal/models"
	"fmt"
)

// Получаем входящие приглашения и уведомления разработчика.
func GetInbox(userID int) (models.InboxResponse, error) {
	response := models.InboxResponse{
		Invites:       []models.InboxInvite{},
		Notifications: []models.InboxNotification{},
	}

	// Показываем ожидающие приглашения, пока набор на слот открыт.
	rows, err := models.UserDB.Query(`
  SELECT i.id, p.id, p.name, s.rolle, s.name, u.Nickname, i.message, i.created_at
  FROM invites i
  JOIN slots s ON s.id = i.slot_id
  JOIN projects p ON p.id = s.project_id
  JOIN users u ON u.id = p.owner_id
  WHERE i.user_id = ? AND i.status = 'pending'
    AND s.status IN ('open', 'review')
    AND p.status != 'finished'
    AND NOT EXISTS (
        SELECT 1 FROM applicationsSlot a
        WHERE a.slot_id = s.id AND a.status = 'approved'
    )
  ORDER BY i.id DESC
 `, userID)
	if err != nil {
		return response, fmt.Errorf("получение приглашений: %w", err)
	}

	for rows.Next() {
		var invite models.InboxInvite
		err = rows.Scan(&invite.ID, &invite.ProjectID, &invite.ProjectName,
			&invite.SlotRole, &invite.SlotName, &invite.OwnerNickname,
			&invite.Message, &invite.CreatedAt)
		if err != nil {
			rows.Close()
			return response, fmt.Errorf("чтение приглашения: %w", err)
		}
		response.Invites = append(response.Invites, invite)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return response, fmt.Errorf("чтение приглашений: %w", err)
	}

	// Получаем только уведомления этого пользователя.
	rows, err = models.UserDB.Query(`
  SELECT n.id, n.text, n.project_id, COALESCE(p.name, ''), n.is_read, n.created_at
  FROM notifications n
  LEFT JOIN projects p ON p.id = n.project_id
  WHERE n.user_id = ?
  ORDER BY n.id DESC
 `, userID)
	if err != nil {
		return response, fmt.Errorf("получение уведомлений: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var notification models.InboxNotification
		err = rows.Scan(&notification.ID, &notification.Text, &notification.ProjectID,
			&notification.ProjectName, &notification.Read, &notification.CreatedAt)
		if err != nil {
			return response, fmt.Errorf("чтение уведомления: %w", err)
		}
		response.Notifications = append(response.Notifications, notification)
	}
	if err = rows.Err(); err != nil {
		return response, fmt.Errorf("чтение уведомлений: %w", err)
	}
	return response, nil
}

// Помечаем уведомления пользователя прочитанными.
func ReadNotifications(userID int) error {
	_, err := models.UserDB.Exec(
		"UPDATE notifications SET is_read = 1 WHERE user_id = ? AND is_read = 0",
		userID,
	)
	if err != nil {
		return fmt.Errorf("прочтение уведомлений: %w", err)
	}
	return nil
}
