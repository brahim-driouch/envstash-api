package interfaces

import (
	"context"

	"github.com/brahim-driouch/envstash.git/internal/models"
)

type UserRepository interface {

	// Read
	FindUserByID(ctx context.Context, id string) (*models.User, error)
	ListUsers(ctx context.Context) ([]*models.User, error)

	// Update
	UpdateUser(ctx context.Context, id string, input *models.UpdateUserInput) (*models.User, error)
	VerifyUser(ctx context.Context, id string) error

	// Delete
	DeleteUser(ctx context.Context, id string) error
}
