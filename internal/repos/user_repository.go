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

func NewUserRepository(db *pgxpool.Pool) (interfaces.UserRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("queriesbase pool cannot be nil")
	}

	// Return concrete implementation as interface
	return &userRepository{
		db: db,
	}, nil
}

func (r *userRepository) CreateUser(ctx context.Context, input *models.CreateUserInput, passwordHash string) (*models.User, error) {
	// Prepare the SQL query for inserting the user
	var u models.User
	err := r.db.QueryRow(
		ctx,
		queries.UserQueries.InsertUser,
		input.Fullname,
		input.Email,
		passwordHash,
	).Scan(
		&u.ID,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(ctx, queries.UserQueries.FindUserByEmail, email).Scan(
		&u.ID,
		&u.Fullname,
		&u.Email,
		&u.PasswordHash,
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
func (r *userRepository) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, queries.UserQueries.UserExists, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	fmt.Printf("%v", exists)
	return exists, nil
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
