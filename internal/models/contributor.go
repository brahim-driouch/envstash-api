package models

import "time"

type Contributor struct {
	ID           string     `json:"id" db:"id"`
	Fullname     string     `json:"fullname"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	IsVerified   bool       `json:"is_verified"`
	UserID       string     `json:"user_id" db:"user_id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`
}

type CreateContributorRequest struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type ContributorWithTeamsAndProjects struct {
	Contributor Contributor
	Teams       []Team    `json:"teams"`
	Projects    []Project `json:"projects"`
}
type ContributorWithTeams struct {
	Contributor Contributor `json:"contributor"`
	Teams       []Team      `json:"teams"`
}

type ContributorWithProjects struct {
	Contributor Contributor `json:"contributor"`
	Projects    []Project   `json:"projects"`
}

type ContributorResponse struct {
	ID         string     `json:"id"`
	Fullname   string     `json:"fullname"`
	Email      string     `json:"email"`
	IsVerified bool       `json:"is_verified"`
	UserID     string     `json:"user_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

func (c *Contributor) ToResponse() ContributorResponse {
	return ContributorResponse{
		ID:         c.ID,
		Fullname:   c.Fullname,
		Email:      c.Email,
		IsVerified: c.IsVerified,
		UserID:     c.UserID,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
		DeletedAt:  c.DeletedAt,
	}
}
