package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRepository struct {
	db *pgxpool.Pool
}

func NewStatsRepository(db *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{
		db: db,
	}
}

func (r *StatsRepository) GetUserStats(ctx context.Context, userID string) (*models.UserStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}
	var stats models.UserStats
	err := r.db.QueryRow(ctx, queries.Stats.AllStats, userID).Scan(&stats.ProjectsCount, &stats.TeamsCount, &stats.MembersCount, &stats.EnvVarsCount)
	if err != nil {
		return nil, err
	}
	if stats == (models.UserStats{}) {
		return nil, errors.New("no stats found for user")
	}
	if stats.ProjectsCount == 0 && stats.TeamsCount == 0 && stats.MembersCount == 0 && stats.EnvVarsCount == 0 {
		return nil, errors.New("no stats found for user")
	}

	return &stats, nil
}
