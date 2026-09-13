package models

import (
	"time"
)

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"target"`
	OwnerID     int       `json:"ownerID"`
	Privacy     string    `json:"privacy"` // либо публичная комната, либо приватная privat либо public
	Status      string    `json:"status"`  //на каком этапе проект (драфт/в разработке/ревью/закрыт/заморожен/приостановлен)
	CreatAt     time.Time `json:"creat_at"`
}

type Slot struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Rolle     string    `json:"rolle"`
	StackID   []int     `json:"stack"`
	Status    string    `json:"status"`
	CreatAt   time.Time `json:"creat_at"`
}

// dto
type DTOProject struct {
	Name    string    `json:"name"`
	Target  string    `json:"target"`
	Privacy string    `json:"privacy"` // public или private
	Status  string    `json:"status"`  // draft, worling, finish
	Slots   []DTOSlot `json:"slots"`
}

type DTOSlot struct {
	Rolle  string `json:"rolle"`
	Stack  []int  `json:"stack"`
	Status string `json:"status"`
}

// dto для ответа
type ProjectSlotOut struct {
	Project Project `json:"project"`
	Slots   []Slot  `json:"slots"`
}

// dto для ответа
type UpdateStatusRoom struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
