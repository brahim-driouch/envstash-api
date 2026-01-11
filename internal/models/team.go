package models

import "time"

type Team struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	UserID      string    `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateTeamRequest struct {
	Name        string `json:"name"`
	UserID      string `json:"user_id"`
	Description string `json:"description,omitempty"`
}
