package repository

import (
	"context"
	"fmt"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/queries"
	"github.com/brahim-driouch/envstash.git/internal/repos/interfaces"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) interfaces.UserRepository {

	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(ctx, queries.UserQueries.FindUserByID, id).Scan(
		&u.ID,
		&u.Email,
		&u.Fullname,
		&u.PasswordHash,
		&u.IsVerified,
		&u.IsAdmin,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	fmt.Println(u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	// Implementation goes here
	return nil, nil
}
func (r *userRepository) UpdateUser(ctx context.Context, id string, input *models.UpdateUserInput) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		ctx,
		queries.UserQueries.UpdateUser,
		id,
		input.Fullname,
		input.IsVerified,
		input.IsAdmin,
	).Scan(
		&u.ID,
		&u.Email,
		&u.Fullname,
		&u.IsVerified,
		&u.IsAdmin,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) VerifyUser(ctx context.Context, id string) error {
	// Implementation goes here
	return nil
}
func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	cmdTag, err := r.db.Exec(ctx, queries.UserQueries.DeleteUser, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
