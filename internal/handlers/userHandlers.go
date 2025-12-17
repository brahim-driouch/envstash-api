package handlers

import (
	"fmt"
	"net/http"

	"github.com/brahim-driouch/envbox.git/internal/auth"
	"github.com/brahim-driouch/envbox.git/internal/models"
	repository "github.com/brahim-driouch/envbox.git/internal/repos"
	"github.com/brahim-driouch/envbox.git/internal/validators"
	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {

	var newUser = models.CreateUserInput{}

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	validationResult := validators.ValidateNewUserFields(newUser)
	errorMessages := make([]string, len(validationResult))
	for i, err := range validationResult {
		errorMessages[i] = err.Error()
	}
	if len(errorMessages) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": errorMessages,
		})
		return
	}

	userRepo := repository.NewUserRepository()

	userExists, err := userRepo.UserExists(c.Request.Context(), newUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check if user exists",
		})
		return
	}
	if userExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user with this email already exists",
		})
		return
	}
	hash, err := auth.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	// Set the hashed password
	newUser.Password = string(hash)

	u, err := userRepo.CreateUser(c.Request.Context(), &newUser, newUser.Password)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registration successful",
		"user":    u.ToResponse(),
	})

}

// delete user handler

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user ID is required",
		})
		return
	}

	userRepo := repository.NewUserRepository()
	u, err := userRepo.FindUserByID(c.Request.Context(), userID)
	if err != nil || u == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	deletionError := userRepo.DeleteUser(c.Request.Context(), userID)
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

func UpdateUser(c *gin.Context) {
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
	userRepo := repository.NewUserRepository()
	u, err := userRepo.FindUserByID(c.Request.Context(), updateData.ID)
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
	updatedUser, err := userRepo.UpdateUser(c.Request.Context(), updateData.ID, &updateData)
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

func LoginUser(c *gin.Context) {

	var userLoginInput models.LoginInput
	//check request payload
	if err := c.ShouldBind(&userLoginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request payload",
		})
		return

	}
	// get the user from db
	userRepo := repository.NewUserRepository()
	user, err := userRepo.FindUserByEmail(c.Request.Context(), userLoginInput.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}
	//if we have the user , compare passwords
	isValidPassword := auth.ComparePasswords(userLoginInput.Password, user.PasswordHash)
	if !isValidPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}
	// genetate token
	var userSub = auth.TokenSub{
		Id:         user.ID,
		Fullname:   user.Fullname,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		IsAdmin:    user.IsAdmin,
	}
	token, tokenErr := auth.GenerateToken(userSub)
	if tokenErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error creating login token, please try again later or report the error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"data": gin.H{

			"token": token,
		},
	})

}
