package main

import (
	"fmt"

	"github.com/brahim-driouch/envbox.git/config"
	"github.com/brahim-driouch/envbox.git/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// Connect to the database
	err := config.ConnectDB()

	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}
	fmt.Println("Connected to the database successfully")

	// Remember to close the database pool when the application exits
	defer config.Pool.Close()
	// Your application logic here
	r := gin.Default()

	//cors
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))
	// Define your routes and handlers here
	apiV1 := r.Group("/api/v1")

	//pulic routes
	apiV1.POST("/users/register", routes.UserRoutes["registerUser"])
	apiV1.POST("/users/login", routes.UserRoutes["loginUser"])
	//protected routes

	apiV1.DELETE("/users/delete/:id", routes.UserRoutes["deleteUser"])
	apiV1.PUT("/users/update/:id", routes.UserRoutes["updateUser"])

	// Start the server
	r.Run()

}
