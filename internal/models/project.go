package models

import (
	"time"
)

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     int       `json:"ownerID"`
	Privacy     string    `json:"privacy"` // public или private
	Status      string    `json:"status"`  // draft, working, finished
	CreatAt     time.Time `json:"creatAt"`
}

type Slot struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"projectID"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Rolle       string    `json:"rolle"`
	StackID     []int     `json:"stackID"`
	UserID      *int      `json:"userID"` // исполнитель, null — никого
	Status      string    `json:"status"` // open, review, close, done
	CreatAt     time.Time `json:"creatAt"`
}

// dto
type DTOProject struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Privacy     string    `json:"privacy"` // public или private
	Status      string    `json:"status"`  // draft, working, finished
	Slots       []DTOSlot `json:"slots"`
}

type DTOSlot struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rolle       string `json:"rolle"`
	StackID     []int  `json:"stackID"`
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
