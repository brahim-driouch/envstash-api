package models

import "time"

type RefreshToken struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expiresAt"`
	CreatedAt time.Time  `json:"createdAt"`
	RevokedAt *time.Time `json:"revokedAt"`
	IPAddress string     `json:"ipAddress"`
	UserAgent string     `json:"userAgent"`
}
type AuthToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
