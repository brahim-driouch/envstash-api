package repository

import (
	"context"
	"errors"
	"log"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	db *pgxpool.Pool // assuming you have a DB type
}

func NewTeamRepository(db *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{db: db}
}

// Implement the TeamRepository interface methods here
func (r *TeamRepository) CreateTeam(ctx context.Context, team *models.CreateTeamRequest) error {
	if team.Name == "" {
		return errors.New("team name is required")
	}
	if team.UserID == "" {
		return errors.New("user id is required")
	}
	cmdTag, err := r.db.Exec(ctx, queries.TeamQueries.CreateTeam, team.Name, team.Description, team.UserID)
	if err != nil {
		log.Println(err)
		return err
	}
	if cmdTag.RowsAffected() == 0 {

		log.Println("No rows affected when creating team")
		log.Printf("Team: %+v", team)
		return errors.New("failed to create team")
	}
	return nil
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, teamID string, userID string) (*models.Team, error) {
	var team models.Team
	err := r.db.QueryRow(ctx, queries.TeamQueries.GetTeamByID, teamID, userID).Scan(&team.ID, &team.Name, &team.Description, &team.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &team, nil
}
func (r *TeamRepository) GetTeamsByUserID(ctx context.Context, userID string) ([]models.Team, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	// 1. Get teams
	teamRows, err := r.db.Query(ctx, queries.TeamQueries.GetTeamsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer teamRows.Close()

	teams := []models.Team{}
	for teamRows.Next() {
		var team models.Team
		err := teamRows.Scan(&team.ID, &team.Name, &team.Description, &team.CreatedAt)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	if err := teamRows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}
func (r *TeamRepository) UpdateTeam(ctx context.Context, team *models.Team) (*models.Team, error) {
	// Implementation goes here
	return nil, nil
}

func (r *TeamRepository) DeleteTeam(ctx context.Context, teamID string) error {
	_, err := r.db.Exec(ctx, queries.TeamQueries.DeleteTeam, teamID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TeamRepository) GetTeamsByUserIDs(ctx context.Context, userIDs []string) ([]models.Team, error) {
	if len(userIDs) == 0 {
		return []models.Team{}, nil
	}

	// 1. Get teams
	teamRows, err := r.db.Query(ctx, queries.TeamQueries.GetTeamsByUserID, userIDs)
	if err != nil {
		return nil, err
	}
	defer teamRows.Close()

	teams := []models.Team{}
	for teamRows.Next() {
		var team models.Team
		err := teamRows.Scan(&team.ID, &team.Name, &team.Description, &team.CreatedAt)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	if err := teamRows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}
