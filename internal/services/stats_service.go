package services

import (
	"context"
	"errors"
	"log"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/repos/interfaces"
)

type StatsService struct {
	statsRepo interfaces.StatsRepository
}

func NewStatsService(repo interfaces.StatsRepository) *StatsService {
	return &StatsService{
		statsRepo: repo,
	}
}

func (s *StatsService) GetUserStats(ctx context.Context, userID string) (*models.UserStats, error) {
	if userID == "" {
		return nil, errors.New("invalid user id")
	}
	stats, err := s.statsRepo.GetUserStats(ctx, userID)
	if err != nil {
		log.Println("Error getting user stats:", err)

		return nil, err
	}

	return stats, nil
}
