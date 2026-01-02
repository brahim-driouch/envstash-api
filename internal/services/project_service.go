package services

import (
	"context"
	"fmt"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/repos/interfaces"
)

// project service errors

var (
	ErrProjectNameExists = fmt.Errorf("project with this name already exists")
	ErrInvalidUserID     = fmt.Errorf("invalid user ID")
	ErrInvalidPayload    = fmt.Errorf("invalid create project payload")
)

type ProjectService struct {
	projectRepo interfaces.ProjectRepository
}

func NewProjectService(projectRepo interfaces.ProjectRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	}
}

// project service

func (s *ProjectService) CreateProject(ctx context.Context, userID string, createProjectPayload *models.CreateProjectRequest) error {
	// validate args
	if userID == "" {
		return ErrInvalidUserID
	}
	if createProjectPayload == nil {
		return ErrInvalidPayload
	}
	// check if user has already a project with the same name
	// find a project with the same name
	project, err := s.projectRepo.GetProjectByName(ctx, createProjectPayload.Name, userID)
	if err == nil && project != nil {
		return ErrProjectNameExists

	}
	err = s.projectRepo.CreateProject(ctx, createProjectPayload, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProjectService) GetProjectsByUserID(ctx context.Context, userID string) (*[]models.Project, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	projects, err := s.projectRepo.GetAllProjects(ctx, userID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}
