package services

import (
	"context"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/repos/interfaces"
)

type DashboardService struct {
	memberRepo  interfaces.ContributorRepository
	teamRepo    interfaces.TeamRepository
	projectRepo interfaces.ProjectRepository
}

func NewDashboardService(contribRepo interfaces.ContributorRepository, teamRepo interfaces.TeamRepository, projectRepo interfaces.ProjectRepository) *DashboardService {
	return &DashboardService{
		memberRepo:  contribRepo,
		teamRepo:    teamRepo,
		projectRepo: projectRepo,
	}
}

func (s *DashboardService) GetUserTeams(ctx context.Context, userID string) ([]models.TeamWithMembersAndProjects, error) {
	return nil, nil
	// if userID == "" {
	// 	return nil, errors.New("user ID is required")
	// }

	// // 1. Get user teams
	// teams, err := s.teamRepo.GetTeamsByUserID(ctx, userID)
	// if err != nil {
	// 	return nil, err
	// }

	// if len(teams) == 0 {
	// 	return []models.TeamWithMembersAndProjects{}, nil
	// }

	// // 2. Extract team IDs
	// teamIDs := make([]string, len(teams))
	// for i, team := range teams {
	// 	teamIDs[i] = team.ID
	// }

	// // 3. Batch fetch members for all teams (1 query instead of N)
	// membersByTeam, err := s.memberRepo.GetMembersByTeamIDs(ctx, teamIDs)
	// if err != nil {
	// 	return nil, err
	// }

	// // 4. Batch fetch projects for all teams (1 query instead of N)
	// //projectsByTeam, err := s.projectRepo.GetProjectsByTeamIDs(ctx, teamIDs)
	// if err != nil {
	// 	return nil, err
	// }

	// // 5. Build response
	// // response := make([]models.TeamWithMembersAndProjects, len(teams))
	// // for i, team := range teams {
	// // 	response[i] = models.TeamWithMembersAndProjects{
	// // 		Team:     team,
	// // 		Members:  membersByTeam[team.ID],  // O(1) lookup
	// // 		Projects: projectsByTeam[team.ID], // O(1) lookup
	// // 	}
	// // }

	// return nil, nil
}
