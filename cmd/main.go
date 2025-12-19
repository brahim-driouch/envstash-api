package main

import (
	"fmt"
	"log"
	"time"

	"github.com/brahim-driouch/envstash.git/config"
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
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Define your routes and handlers here
	apiV1 := r.Group("/api/v1")
	//instantiate user service and user repo
	userRepository, err := repository.NewUserRepository(dbPool)
	if err != nil {
		fmt.Println("Could not instantiate user repository %w", err)
	}
	userService := services.NewUserHandler(userRepository)
	//pulic routes
	apiV1.POST("/users/register", userService.RegisterUser)
	apiV1.POST("/users/login", userService.LoginUser)
	//protected routes

	apiV1.DELETE("/users/delete/:id", userService.DeleteUser)
	apiV1.PUT("/users/update/:id", userService.UpdateUser)

	// Start the server
	r.Run()

}
