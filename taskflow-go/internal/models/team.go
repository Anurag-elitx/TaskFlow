package models

import (
	"time"
	"github.com/google/uuid"
)

type Team struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type TeamMember struct {
	TeamID   uuid.UUID `json:"team_id"`
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}
