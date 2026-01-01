package handlers

import (
	"log"
	"net/http"

	"github.com/brahim-driouch/envstash.git/internal/services"
	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService *services.StatsService
}

func NewStatsHandler(service *services.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: service,
	}
}

func (h *StatsHandler) GetUserStats(c *gin.Context) {
	userID := c.Param("id")
	log.Println(userID)

	userContextID, exists := c.Get("userId")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userContextID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// TODO: Use userID and userContextID to get user stats
	ctx := c.Request.Context()
	stats, err := h.statsService.GetUserStats(ctx, userID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user stats"})
		return
	}
	if stats == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User stats not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}
