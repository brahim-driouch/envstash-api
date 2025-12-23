package handlers

import (
	"net/http"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: s,
	}
}
func (h *AuthHandler) CreateRefreshToken(c *gin.Context) {
	var refreshToken models.RefreshToken
	if err := c.ShouldBindJSON(&refreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	ctx := c.Request.Context()
	err := h.authService.CreateRefreshToken(ctx, &refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while creating refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Refresh token created successfully", "data": refreshToken})
}
