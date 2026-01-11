package interfaces

import (
	"context"

	"github.com/brahim-driouch/envstash.git/internal/models"
)

type TeamRepository interface {
	// Team operations
	CreateTeam(ctx context.Context, team *models.CreateTeamRequest) error
	GetTeamByID(ctx context.Context, id string) (*models.Team, error)
	GetTeamsByUserID(ctx context.Context, userID string) ([]models.Team, error)
	UpdateTeam(ctx context.Context, team *models.Team) (*models.Team, error)
	DeleteTeam(ctx context.Context, teamID string) error
}
