package repository

import (
	"context"
	"errors"
	"log"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, project *models.CreateProjectRequest, userID string) error {
	_, err := r.db.Exec(ctx, queries.ProjectQueries.CreateProject, project.Name, project.Description, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProjectRepository) GetProjectByID(ctx context.Context, id string) (*models.Project, error) {
	return nil, nil
}

// TODO CREATE INDEXES ON MEMBERS TABLE
func (r *ProjectRepository) GetAllProjects(ctx context.Context, userID string) (*[]models.Project, error) {
	var projects []models.Project
	rows, err := r.db.Query(ctx, queries.ProjectQueries.GetAllProjects, userID)
	if err != nil {
		log.Println(err)

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var project models.Project
		err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.UserID, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return &projects, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, id string, project *models.UpdateProjectRequest) error {

	return nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id string) error {

	return nil
}

func (r *ProjectRepository) GetProjectsByUserID(ctx context.Context, userID string) ([]models.Project, error) {
	// TODO: Implement this method
	return nil, nil
}

func (r *ProjectRepository) GetProjectByName(ctx context.Context, name string, userID string) (*models.Project, error) {
	return nil, nil
}

func (r *ProjectRepository) GetProjectMembers(ctx context.Context, projectID string) (*[]models.ProjectMember, error) {
	if projectID == "" {
		return nil, errors.New("invalid prject ID")
	}
	var projectMembers []models.ProjectMember
	rows, err := r.db.Query(ctx, queries.ProjectQueries.GetProjectMembers, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var projectMember models.ProjectMember
		err := rows.Scan(&projectMember.ID, &projectMember.ProjectID, &projectMember.UserID, &projectMember.CreatedAt)
		if err != nil {
			return nil, err
		}
		projectMembers = append(projectMembers, projectMember)
	}
	return &projectMembers, nil
}

func (r *ProjectRepository) GetProjectVars(ctx context.Context, projectID string) (*[]models.EnvironmentVariable, error) {
	if projectID == "" {
		return nil, errors.New("invalid project ID")
	}
	rows, err := r.db.Query(ctx, queries.ProjectQueries.GetProjectVars, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projectVars []models.EnvironmentVariable

	for rows.Next() {
		var projectVar models.EnvironmentVariable
		err := rows.Scan(&projectVar.ID, &projectVar.ProjectID, &projectVar.Key, &projectVar.Value, &projectVar.CreatedAt, &projectVar.UpdatedAt, &projectVar.CreatedBy)
		if err != nil {
			return nil, err
		}
		projectVars = append(projectVars, projectVar)
	}
	return &projectVars, nil
}
