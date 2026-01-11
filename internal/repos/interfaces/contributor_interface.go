package interfaces

import (
	"context"

	"github.com/brahim-driouch/envstash.git/internal/models"
)

type ContributorRepository interface {
	// Add member-related methods here
	GetContributorsByTeamIDs(ctx context.Context, teamIDs []string) (map[string][]models.Contributor, error)
	CreateContributor(ctx context.Context, contributor models.CreateContributorRequest) error
	GetContributorsByUserID(ctx context.Context, userID string) ([]models.ContributorResponse, error)
	GetContributorByID(ctx context.Context, contributorID string) (*models.ContributorResponse, error)
}
