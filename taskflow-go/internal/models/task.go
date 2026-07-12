package models

import (
	"time"
	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	TeamID      uuid.UUID  `json:"team_id"`
	DueDate     *time.Time `json:"due_date"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskFilters struct {
	Status     *string
	Priority   *string
	AssigneeID *uuid.UUID
	Page       int
	Limit      int
}
