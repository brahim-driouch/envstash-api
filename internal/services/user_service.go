package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brahim-driouch/envstash.git/internal/auth"
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

func (s *UserService) RegisterUser(ctx context.Context, input *models.CreateUserInput) (*models.User, error) {
	validationError := validators.ValidateNewUserFields(*input)
	if validationError != nil {
		return nil, validationError
	}

	userExists, err := s.userRepo.UserExists(ctx, input.Email)
	if err != nil {
		return nil, ErrUnexpected
	}
	if userExists {
		return nil, ErrUserExists
	}
	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, ErrUnexpected
	}
	// Set the hashed password
	input.Password = string(hash)

	u, err := s.userRepo.CreateUser(ctx, input, input.Password)

	if err != nil {
		return nil, ErrUnexpected
	}
	return u, nil

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

func (s *UserService) LoginUser(ctx context.Context, userLoginInput models.LoginInput) (*models.AuthToken, error) {
	// get the user from db
	user, err := s.userRepo.FindUserByEmail(ctx, userLoginInput.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	//if we have the user , compare passwords
	isValidPassword := auth.ComparePasswords(userLoginInput.Password, user.PasswordHash)
	if !isValidPassword {
		return nil, ErrInvalidCredentials
	}
	// genetate token
	var userSub = auth.TokenSub{
		Id:         user.ID,
		Fullname:   user.Fullname,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		IsAdmin:    user.IsAdmin,
	}
	// set access token err to 15 minutes
	accessToken, accessTokenErr := auth.GenerateToken(userSub, 15)
	//set the refressh token for 30 dayas
	refreshToken, refreshTokenErr := auth.GenerateToken(userSub, 60*24*15)

	if accessTokenErr != nil || refreshTokenErr != nil {
		return nil, ErrUnexpected
	}
	return &models.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
