package interfaces

import (
	"context"

	"github.com/brahim-driouch/envstash.git/internal/models"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, project *models.CreateProjectRequest, userID string) error
	GetProjectByID(ctx context.Context, id string) (*models.Project, error)
	GetProjectByName(ctx context.Context, name string, userID string) (*models.Project, error)
	GetAllProjects(ctx context.Context, userID string) (*[]models.Project, error)
	UpdateProject(ctx context.Context, id string, project *models.UpdateProjectRequest) error
	DeleteProject(ctx context.Context, id string) error
}
