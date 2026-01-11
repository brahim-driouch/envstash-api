package handlers

import (
	"net/http"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/services"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(teamService *services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

// handler functions

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("userId")
	if !exists {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var newTeam models.CreateTeamRequest
	if err := c.ShouldBindJSON(&newTeam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	err := h.teamService.CreateTeam(ctx, userID.(string), &newTeam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team"})
		return
	}

	c.JSON(http.StatusCreated, newTeam)
}

func (h *TeamHandler) GetTeamsByUserID(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	reqUserID := c.Query("user_id")
	if reqUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter is required"})
		return
	}
	if reqUserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot access teams of another user"})
		return
	}

	teams, err := h.teamService.GetTeamsByUserID(ctx, reqUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get teams"})
		return
	}

	c.JSON(http.StatusOK, teams)
}

// delete team by id

func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	ctx := c.Request.Context()
	teamID := c.Param("team_id")
	if teamID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team_id is required"})
		return
	}
	err := h.teamService.DeleteTeam(ctx, teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete team"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}
