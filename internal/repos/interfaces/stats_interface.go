package interfaces

import (
	"context"

	"github.com/brahim-driouch/envstash.git/internal/models"
)

type StatsRepository interface {
	GetUserStats(ctx context.Context, userID string) (*models.UserStats, error)
}
