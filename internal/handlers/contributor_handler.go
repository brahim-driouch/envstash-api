package handlers

import (
	"net/http"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/services"
	"github.com/gin-gonic/gin"
)

type ContributorHandler struct {
	contributorService *services.ContributorService
}

func NewContributorHandler(contributorService *services.ContributorService) *ContributorHandler {
	return &ContributorHandler{
		contributorService: contributorService,
	}
}

func (h *ContributorHandler) GetContributorByID(c *gin.Context) {
	contributorID := c.Param("id")
	if contributorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contributor ID is required"})
		return
	}
	ctx := c.Request.Context()
	contributor, err := h.contributorService.GetContributorByID(ctx, contributorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"contributor": contributor})
}
func (h *ContributorHandler) CreateContributor(c *gin.Context) {
	var contributor models.CreateContributorRequest
	if err := c.ShouldBindJSON(&contributor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx := c.Request.Context()
	contributor.UserID = userId.(string)
	err := h.contributorService.CreateContributor(ctx, contributor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "contributor created successfully"})
}

func (h *ContributorHandler) GetContributorsByUserID(c *gin.Context) {
	ctx := c.Request.Context()
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	reqUserID := c.Query("user_id")
	if reqUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	if userId != reqUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	contributors, err := h.contributorService.GetContributorsByUserID(ctx, reqUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"contributors": contributors})
}
