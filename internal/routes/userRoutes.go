package routes

import (
	"github.com/brahim-driouch/envbox.git/internal/handlers"
	"github.com/gin-gonic/gin"
)

var UserRoutes = map[string]gin.HandlerFunc{
	"registerUser": handlers.RegisterUser,
	"deleteUser":   handlers.DeleteUser,
	"updateUser":   handlers.UpdateUser,
	"loginUser":    handlers.LoginUser,
}
