package repository

import (
	"context"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

// Create refresh token for user and a device
func (r *AuthRepository) CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	// Use QueryRow with RETURNING to get the generated ID
	err := r.db.QueryRow(
		ctx,
		queries.AuthQueries.CreateRefreshToken,
		refreshToken.UserID,
		refreshToken.Token,
		refreshToken.ExpiresAt,
		refreshToken.CreatedAt,
		refreshToken.IPAddress,
		refreshToken.UserAgent,
	).Scan(&refreshToken.ID)

	return err
}

// Find a valid (non-revoked, non-expired) refresh token
func (r *AuthRepository) FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken

	err := r.db.QueryRow(
		ctx,
		queries.AuthQueries.FindRefreshToken,
		token,
	).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
		&refreshToken.RevokedAt,
		&refreshToken.IPAddress,
		&refreshToken.UserAgent,
	)

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

// Revoke a specific refresh token (logout)
// BUG FIX: Use Exec, not QueryRow (UPDATE doesn't return rows)
func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.Exec(ctx, queries.AuthQueries.RevokeRefreshToken, token)
	return err
}

// Revoke all tokens for a user (logout from all devices)
// BUG FIX: Use Exec, not QueryRow
func (r *AuthRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, queries.AuthQueries.RevokeAllUserTokens, userID)
	return err
}

// Delete expired and revoked tokens (cleanup job)
// BUG FIX: Use Exec, not QueryRow
func (r *AuthRepository) DeleteExpiredTokens(ctx context.Context) error {
	_, err := r.db.Exec(ctx, queries.AuthQueries.DeleteExpiredTokens)
	return err
}

// Find all active sessions for a user
// BUG FIX: Use Query (not QueryRow) for multiple rows
func (r *AuthRepository) FindActiveUserTokens(ctx context.Context, userID string) ([]models.RefreshToken, error) {
	rows, err := r.db.Query(ctx, queries.AuthQueries.FindActiveUserTokens, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var refreshTokens []models.RefreshToken

	for rows.Next() {
		var token models.RefreshToken
		err := rows.Scan(
			&token.ID,
			&token.Token,
			&token.ExpiresAt,
			&token.CreatedAt,
			&token.IPAddress,
			&token.UserAgent,
		)
		if err != nil {
			return nil, err
		}
		refreshTokens = append(refreshTokens, token)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return refreshTokens, nil
}

// Delete a specific token by ID (for session management UI)
// BUG FIX: Use Exec, not QueryRow
func (r *AuthRepository) DeleteUserToken(ctx context.Context, tokenID string, userID string) error {
	_, err := r.db.Exec(ctx, queries.AuthQueries.DeleteUserToken, tokenID, userID)
	return err
}
