package models

import "time"

type Team struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	UserID      string    `json:"userId" db:"user_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type CreateTeamRequest struct {
	Name        string `json:"name"`
	UserID      string `json:"userId"`
	Description string `json:"description,omitempty"`
}
