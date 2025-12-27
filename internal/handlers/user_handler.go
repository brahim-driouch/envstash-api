package handlers

import (
	"github.com/brahim-driouch/envstash.git/internal/services"
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

func (h *UserHandler) DeleteUser(c *gin.Context) {}
func (h *UserHandler) UpdateUser(c *gin.Context) {}
