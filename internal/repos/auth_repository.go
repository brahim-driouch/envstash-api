package repository

import (
	"context"
	"errors"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/queries"
	"github.com/brahim-driouch/envstash.git/internal/repos/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) interfaces.AuthRepository {
	return &authRepository{
		db: db,
	}
}

//LoginUser

// Create refresh token for user and a device
func (r *authRepository) CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error {
	// Use QueryRow with RETURNING to get the generated ID
	_, err := r.db.Exec(
		ctx,
		queries.AuthQueries.CreateRefreshToken,
		refreshToken.UserID,
		refreshToken.Token,
		refreshToken.ExpiresAt,
		refreshToken.CreatedAt,
		refreshToken.IPAddress,
		refreshToken.UserAgent,
	)

	return err
}

// Find a valid (non-revoked, non-expired) refresh token
func (r *authRepository) FindRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
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
func (r *authRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.Exec(ctx, queries.AuthQueries.RevokeRefreshToken, token)
	return err
}

// Revoke all tokens for a user (logout from all devices)
// BUG FIX: Use Exec, not QueryRow
func (r *authRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, queries.AuthQueries.RevokeAllUserTokens, userID)
	return err
}

// Delete expired and revoked tokens (cleanup job)
// BUG FIX: Use Exec, not QueryRow
func (r *authRepository) DeleteExpiredTokens(ctx context.Context) error {
	_, err := r.db.Exec(ctx, queries.AuthQueries.DeleteExpiredTokens)
	return err
}

// Find all active sessions for a user
// BUG FIX: Use Query (not QueryRow) for multiple rows
func (r *authRepository) FindActiveUserTokens(ctx context.Context, userID string) (*[]models.RefreshToken, error) {
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

	return &refreshTokens, nil
}

// Delete a specific token by ID (for session management UI)
// BUG FIX: Use Exec, not QueryRow
func (r *authRepository) DeleteUserToken(ctx context.Context, tokenID string, userID string) error {
	_, err := r.db.Exec(ctx, queries.AuthQueries.DeleteUserToken, tokenID, userID)
	return err
}

func (r *authRepository) CreateUser(ctx context.Context, input *models.CreateUserInput, passwordHash string) (*models.User, error) {
	return nil, nil
}
func (r *authRepository) UserExists(ctx context.Context, email string) (bool, error) {
	return false, nil
}
func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := models.User{}
	err := r.db.QueryRow(ctx, queries.AuthQueries.FindUserByEmail, email).Scan(&user.ID, &user.Fullname, &user.Email, &user.PasswordHash, &user.IsVerified, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	if userID == "" {
		return nil, errors.New("no userID provided")
	}
	var user models.User
	err := r.db.QueryRow(ctx, queries.AuthQueries.FindUserByID, userID).Scan(&user.ID, &user.Fullname, &user.Email, &user.PasswordHash, &user.IsVerified, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
