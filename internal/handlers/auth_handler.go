package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/brahim-driouch/envstash.git/internal/auth"
	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/services"
	"github.com/brahim-driouch/envstash.git/internal/validators"
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
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var newUser = models.CreateUserInput{}

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	user, err := h.authService.RegisterUser(c, &newUser)

	if err != nil {
		switch err {
		case services.ErrInvalidCredentials:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		case services.ErrUnexpected:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		case services.ErrUserExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		case validators.ErrPasswordLength,
			validators.ErrInvalidEmail,
			validators.ErrMissingFields,
			validators.ErrNameStringLength,
			validators.ErrPasswordLength,
			validators.ErrPasswordMatch:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occured, please try again later or report the error"})

		}
		return

	}
	c.JSON(http.StatusCreated, gin.H{"message": "registred successfully", "data": gin.H{"user": user.ToResponse()}})
}
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var loginInput models.LoginInput
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	tokens, err := h.authService.LoginUser(c.Request.Context(), loginInput)
	fmt.Println(err)
	if err != nil {
		switch err {

		case services.ErrInvalidCredentials,
			services.ErrUserNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		case services.ErrUnexpected:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occured, please try again later or report the error"})
			return
		}

	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully", "data": tokens})
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

func (h *AuthHandler) LogoutUser(c *gin.Context) {
	ctx := c.Request.Context()
	token := auth.GetAuthorizationHeader(c)
	log.Println(token)

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing refresh token"})
		return
	}
	err := h.authService.LogoutUser(ctx, token)
	if err != nil {
		log.Println(err, token)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while logging out"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// retrieve session
func (h *AuthHandler) GetSession(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return

	}
	newAccessToken := c.GetHeader("X-New-Access-Token")
	if newAccessToken != "" {
		c.Header("X-New-Access-Token", newAccessToken)
	}
	c.JSON(http.StatusOK, gin.H{"message": "User session retrieved successfully", "user": user})
}
