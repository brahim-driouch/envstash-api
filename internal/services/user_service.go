package services

import (
	"fmt"
	"net/http"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/repos/interfaces"
	"github.com/brahim-driouch/envstash.git/internal/validators"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// delete user handler

func (s *UserService) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user ID is required",
		})
		return
	}

	u, err := s.userRepo.FindUserByID(c.Request.Context(), userID)
	if err != nil || u == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	deletionError := s.userRepo.DeleteUser(c.Request.Context(), userID)
	if deletionError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func (s *UserService) UpdateUser(c *gin.Context) {
	var updateData models.UpdateUserInput
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	err := validators.ValidateUpdateUserFields(updateData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	//checl if user exists
	u, err := s.userRepo.FindUserByID(c.Request.Context(), updateData.ID)
	fmt.Println(updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check if user exists",
		})
		return
	}
	if u.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user with this email does not exist",
		})
		return
	}
	// update user
	updatedUser, err := s.userRepo.UpdateUser(c.Request.Context(), updateData.ID, &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    updatedUser.ToResponse(),
	})

}
