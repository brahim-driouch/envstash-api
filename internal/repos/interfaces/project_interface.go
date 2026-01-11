package interfaces

import (
	"context"

	"github.com/brahim-driouch/envstash.git/internal/models"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, project *models.CreateProjectRequest, userID string) error
	GetProjectByID(ctx context.Context, projectId string) (*models.Project, error)
	GetProjectByName(ctx context.Context, name string, userID string) (*models.Project, error)
	GetAllProjects(ctx context.Context, userID string) (*[]models.Project, error)
	UpdateProject(ctx context.Context, projectId string, project *models.UpdateProjectRequest) error
	DeleteProject(ctx context.Context, projectId string) error
	GetProjectMembers(ctx context.Context, projectId string) (*[]models.ProjectMember, error)
	GetProjectVars(ctx context.Context, projectId string) (*[]models.EnvironmentVariable, error)
}
