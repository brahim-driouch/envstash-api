package main

import (
	"fmt"
	"log"
	"time"

	"github.com/brahim-driouch/envstash.git/config"
	"github.com/brahim-driouch/envstash.git/internal/auth"
	"github.com/brahim-driouch/envstash.git/internal/handlers"
	repository "github.com/brahim-driouch/envstash.git/internal/repos"
	"github.com/brahim-driouch/envstash.git/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	// Connect to the database
	dbPool, err := config.ConnectDB()

	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return
	}
	fmt.Println("Connected to the database successfully")

	// Remember to close the database pool when the application exits
	defer dbPool.Close()
	// Your application logic here
	r := gin.Default()

	//cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-New-Access-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Define your routes and handlers here
	apiV1 := r.Group("/api/v1")
	//instantiate user service and user repo
	userRepository := repository.NewUserRepository(dbPool)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	authRepository := repository.NewAuthRepository(dbPool)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService)
	//pulic routes
	apiV1.POST("/auth/register", authHandler.RegisterUser)
	apiV1.POST("/auth/login", authHandler.LoginUser)
	apiV1.POST("/auth/logout", authHandler.LogoutUser)
	//get current session
	apiV1.GET("/auth/session", auth.AuthMiddleware(authService), authHandler.GetSession)
	//protected routes

	//delete user
	apiV1.DELETE("/users/delete/:id", auth.AuthMiddleware(authService), userHandler.DeleteUser)
	//update user
	apiV1.PUT("/users/update/:id", auth.AuthMiddleware(authService), userHandler.UpdateUser)

	// Start the server
	r.Run()

}
