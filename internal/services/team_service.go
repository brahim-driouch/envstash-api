package services

import (
	"context"
	"errors"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/repos/interfaces"
)

type TeamService struct {
	teamRepo interfaces.TeamRepository
}

func NewTeamService(teamRepo interfaces.TeamRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
	}
}

// service

func (s *TeamService) CreateTeam(ctx context.Context, userID string, team *models.CreateTeamRequest) error {
	if userID == "" {
		return errors.New("user ID is required")
	}
	if len(team.Name) < 2 {
		return errors.New("team name must be at least 2 characters long")
	}
	team.UserID = userID
	err := s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		return err
	}
	return nil
}

func (s *TeamService) GetTeamsByUserID(ctx context.Context, userID string) ([]models.Team, error) {

	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	teams, err := s.teamRepo.GetTeamsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

// delete team

func (s *TeamService) DeleteTeam(ctx context.Context, teamID string) error {
	if teamID == "" {
		return errors.New("team ID is required")
	}
	err := s.teamRepo.DeleteTeam(ctx, teamID)
	if err != nil {
		return err
	}
	return nil
}
