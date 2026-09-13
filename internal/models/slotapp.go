package models


import (
	"time"
)


type AplicationSlot struct {
	ID int `json:"id"`
	SlotID int `json:"slot_id"`
	UserID int `json:"user_id"`
	Status string `json:"status"`
	Creat_add  time.Time `json:"creat_add"`
}


type ApplicationResponse struct {
	Message string `json:"message"`
}