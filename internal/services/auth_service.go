package services

import (
	"context"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/repos/interfaces"
)

type AuthService struct {
	authRepo interfaces.AuthRepository
}

func NewAuthService(r interfaces.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: r,
	}
}

func (s *AuthService) CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	return s.authRepo.CreateRefreshToken(ctx, refreshToken)
}

func (s *AuthService) FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	return s.authRepo.FindRefreshToken(ctx, token)
}

func (s *AuthService) RevokeRefreshToken(ctx context.Context, token string) error {
	return s.authRepo.RevokeRefreshToken(ctx, token)
}

func (s *AuthService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return s.authRepo.RevokeAllUserTokens(ctx, userID)
}

func (s *AuthService) DeleteExpiredTokens(ctx context.Context) error {
	return s.authRepo.DeleteExpiredTokens(ctx)
}

func (s *AuthService) FindActiveUserTokens(ctx context.Context, userID string) ([]*models.RefreshToken, error) {
	return s.authRepo.FindActiveUserTokens(ctx, userID)
}

func (s *AuthService) DeleteUserToken(ctx context.Context, tokenID string, userID string) error {
	return s.authRepo.DeleteUserToken(ctx, tokenID, userID)
}
