package repository

import (
	"context"
	"fmt"

	"github.com/brahim-driouch/envbox.git/config"
	"github.com/brahim-driouch/envbox.git/internal/data"
	"github.com/brahim-driouch/envbox.git/internal/models"
)

type UserRepositoryImpl struct {
	BaseRepository
}

func NewUserRepository() UserRepository {
	pool := config.Pool
	if pool == nil {
		panic("Database pool is not initialized")
	}
	return &UserRepositoryImpl{
		BaseRepository: BaseRepository{
			Pool: pool,
		},
	}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, input *models.CreateUserInput, passwordHash string) (*models.User, error) {
	// Prepare the SQL query for inserting the user
	var u models.User
	err := r.Pool.QueryRow(
		ctx,
		data.UserQueries.InsertUser,
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

func (r *UserRepositoryImpl) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.Pool.QueryRow(ctx, data.UserQueries.FindUserByEmail, email).Scan(
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
func (r *UserRepositoryImpl) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.Pool.QueryRow(ctx, data.UserQueries.UserExists, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	fmt.Printf("%v", exists)
	return exists, nil
}

func (r *UserRepositoryImpl) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := r.Pool.QueryRow(ctx, data.UserQueries.FindUserByID, id).Scan(
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

func (r *UserRepositoryImpl) ListUsers(ctx context.Context) ([]*models.User, error) {
	// Implementation goes here
	return nil, nil
}
func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, id string, input *models.UpdateUserInput) (*models.User, error) {
	var u models.User
	err := r.Pool.QueryRow(
		ctx,
		data.UserQueries.UpdateUser,
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

func (r *UserRepositoryImpl) VerifyUser(ctx context.Context, id string) error {
	// Implementation goes here
	return nil
}
func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, id string) error {
	cmdTag, err := r.Pool.Exec(ctx, data.UserQueries.DeleteUser, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
