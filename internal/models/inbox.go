package models

import "time"

// Приглашение, которое ожидает ответа разработчика.
type InboxInvite struct {
	ID            int       `json:"id"`
	ProjectID     int       `json:"projectID"`
	ProjectName   string    `json:"projectName"`
	SlotRole      string    `json:"slotRole"`
	SlotName      string    `json:"slotName"`
	OwnerNickname string    `json:"ownerNickname"`
	Message       string    `json:"message"`
	CreatedAt     time.Time `json:"createdAt"`
}

// Уведомление во входящих пользователя.
type InboxNotification struct {
	ID          int       `json:"id"`
	Text        string    `json:"text"`
	ProjectID   *int      `json:"projectID"`
	ProjectName string    `json:"projectName"`
	Read        bool      `json:"read"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Приглашения и уведомления текущего пользователя.
type InboxResponse struct {
	Invites       []InboxInvite       `json:"invites"`
	Notifications []InboxNotification `json:"notifications"`
}
