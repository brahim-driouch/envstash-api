package handlers

import (
	"log"
	"net/http"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/services"
	"github.com/brahim-driouch/envstash.git/internal/utils"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService *services.ProjectService
}

func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var createProjectPayload models.CreateProjectRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&createProjectPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	userSub, exists := c.Get("user")
	// check the user in the context
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access - user not found in context"})
		return
	}
	user := userSub.(utils.TokenSub)
	err := h.projectService.CreateProject(ctx, user.Id, &createProjectPayload)
	if err != nil {
		log.Println("Error creating project:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "project created successfully", "user": user})
}

func (h *ProjectHandler) GetProjectsByUserID(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	userSub, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access - user not found in context"})
		return
	}
	user := userSub.(utils.TokenSub)
	if user.Id != userId {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: cannot access projects of another user"})
		return
	}
	projects, err := h.projectService.GetProjectsByUserID(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get projects"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}
