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
	authRoutes := apiV1.Group("/auth")
	projectRoutes := apiV1.Group("/projects")
	//instantiate user service and user repo
	userRepository := repository.NewUserRepository(dbPool)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	authRepository := repository.NewAuthRepository(dbPool)
	projectRepository := repository.NewProjectRepository(dbPool)
	authService := services.NewAuthService(authRepository)
	projectService := services.NewProjectService(projectRepository)
	projectHandler := handlers.NewProjectHandler(projectService)
	authHandler := handlers.NewAuthHandler(authService)
	//pulic routes
	authRoutes.POST("/register", authHandler.RegisterUser)
	authRoutes.GET("/verify", authHandler.VerifyEmail)
	authRoutes.POST("/resend-verification", authHandler.ResendVerificationEmail)
	authRoutes.POST("/login", authHandler.LoginUser)
	authRoutes.POST("/logout", authHandler.LogoutUser)

	//protected routes
	{
		// getr current session
		authRoutes.GET("/session", auth.AuthMiddleware(authService), authHandler.GetSession)
		//delete user
		apiV1.DELETE("/users/delete/:id", auth.AuthMiddleware(authService), userHandler.DeleteUser)
		//update user
		apiV1.PUT("/users/update/:id", auth.AuthMiddleware(authService), userHandler.UpdateUser)

		//projects routes
		projectRoutes.POST("/", auth.AuthMiddleware(authService), projectHandler.CreateProject)
		// projectRoutes.GET("/:id", auth.AuthMiddleware(authService), projectHandler.GetProjectByID)
		// projectRoutes.PUT("/:id", auth.AuthMiddleware(authService), projectHandler.UpdateProject)
		// projectRoutes.DELETE("/:id", auth.AuthMiddleware(authService), projectHandler.DeleteProject)
	}

	// Start the server
	r.Run()

}
