package interfaces

import (
	"context"

	"github.com/brahim-driouch/envstash.git/internal/models"
)

type AuthRepository interface {
	// Create
	CreateUser(ctx context.Context, input *models.CreateUserInput, passwordHash string) (*models.User, error)
	UserExists(ctx context.Context, email string) (bool, error)
	CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error
	FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	FindUserByID(ctx context.Context, userID string) (*models.User, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
	DeleteExpiredTokens(ctx context.Context) error
	FindActiveUserTokens(ctx context.Context, userID string) (*[]models.RefreshToken, error)
	DeleteUserToken(ctx context.Context, tokenID string, userID string) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
}
