package handlers

import (
	"fmt"
	"net/http"

	"github.com/brahim-driouch/envstash.git/internal/models"
	"github.com/brahim-driouch/envstash.git/internal/services"
	"github.com/brahim-driouch/envstash.git/internal/validators"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{
		userService: s,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var newUser = models.CreateUserInput{}

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}
	user, err := h.userService.RegisterUser(c, &newUser)

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

func (h *UserHandler) LoginUser(c *gin.Context) {
	var loginInput models.LoginInput
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	tokens, err := h.userService.LoginUser(c.Request.Context(), loginInput)
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
func (h *UserHandler) DeleteUser(c *gin.Context) {}
func (h *UserHandler) UpdateUser(c *gin.Context) {}
func (h *UserHandler) GetSession(c *gin.Context) {
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
