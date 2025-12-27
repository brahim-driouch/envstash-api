// internal/models/user.go
package models

import "time"

// User represents a user in the system
type User struct {
	ID           string    `json:"id" db:"id"`
	Fullname     string    `json:"fullname" db:"fullname"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never return in JSON
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	IsAdmin      bool      `json:"is_admin" db:"is_admin"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateUserInput is used when creating a user
type CreateUserInput struct {
	Fullname        string `json:"fullname" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// LoginInput is used for login
type LoginInput struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
	UserAgent string `json:"user_agent" binding:"required"`
}

// UpdateUserInput is used when updating a user
type UpdateUserInput struct {
	ID         string  `json:"id" binding:"required"`
	Fullname   *string `json:"fullname,omitempty"`
	IsVerified *bool   `json:"is_verified,omitempty"`
	IsAdmin    *bool   `json:"is_admin,omitempty"`
}

// UserResponse is what gets returned to the client (no password hash)
type UserResponse struct {
	ID         string    `json:"id"`
	Fullname   string    `json:"fullname"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	IsAdmin    bool      `json:"is_admin"`
	CreatedAt  time.Time `json:"created_at"`
}

// ToResponse converts User to UserResponse (strips sensitive data)
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Fullname:   u.Fullname,
		Email:      u.Email,
		IsVerified: u.IsVerified,
		IsAdmin:    u.IsAdmin,
		CreatedAt:  u.CreatedAt,
	}
}
