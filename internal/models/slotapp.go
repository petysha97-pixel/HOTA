package models

import (
	"time"
)

type AplicationSlot struct {
	ID      int       `json:"id"`
	SlotID  int       `json:"slotID"`
	UserID  int       `json:"userID"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
	CreatAt time.Time `json:"creatAt"`
}

type ApplicationResponse struct {
	Message string `json:"message"`
}
