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
	userRoutes := apiV1.Group("/users")
	teamRoutes := apiV1.Group("/teams")
	contributorRoutes := apiV1.Group("/contributors")

	// groupes routes where multiple are needed in the same request
	// dashboardRoutes := apiV1.Group("/dashboard")

	//contributors service
	contributorRepository := repository.NewContributorRepository(dbPool)
	contributorService := services.NewContributorService(contributorRepository)
	contributorHandler := handlers.NewContributorHandler(contributorService)
	//user service
	userRepository := repository.NewUserRepository(dbPool)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)
	//auth services
	authRepository := repository.NewAuthRepository(dbPool)
	authService := services.NewAuthService(authRepository)
	authHandler := handlers.NewAuthHandler(authService)
	//project services
	projectRepository := repository.NewProjectRepository(dbPool)
	projectService := services.NewProjectService(projectRepository)
	projectHandler := handlers.NewProjectHandler(projectService)

	//stats services
	statsRepository := repository.NewStatsRepository(dbPool)
	statsService := services.NewStatsService(statsRepository)
	statsHandler := handlers.NewStatsHandler(statsService)
	// team services
	teamRepository := repository.NewTeamRepository(dbPool)
	teamService := services.NewTeamService(teamRepository)
	teamHandler := handlers.NewTeamHandler(teamService)
	// member services
	// memberRepository := repository.NewMemberRepository(dbPool)
	// memberService := services.NewMemberService(memberRepository)
	// memberHandler := handlers.NewMemberHandler(memberService)
	//aggregator service

	// dashboardService := services.NewDashboardService(memberRepository, teamRepository, projectRepository)
	//public routes
	authRoutes.POST("/register", authHandler.RegisterUser)
	authRoutes.GET("/verify", authHandler.VerifyEmail)
	authRoutes.POST("/resend-verification", authHandler.ResendVerificationEmail)
	authRoutes.POST("/login", authHandler.LoginUser)
	authRoutes.POST("/logout", authHandler.LogoutUser)

	//protected routes
	{
		// get current session
		authRoutes.GET("/session", auth.AuthMiddleware(authService), authHandler.GetSession)
		//delete user
		userRoutes.DELETE("/:id", auth.AuthMiddleware(authService), userHandler.DeleteUser)
		//update user
		userRoutes.PUT("/:id", auth.AuthMiddleware(authService), userHandler.UpdateUser)

		//projects routes
		projectRoutes.POST("/", auth.AuthMiddleware(authService), projectHandler.CreateProject)
		projectRoutes.GET("/", auth.AuthMiddleware(authService), projectHandler.GetProjectsByUserID)
		// projectRoutes.GET("/:id", auth.AuthMiddleware(authService), projectHandler.GetProjectByID)
		// projectRoutes.PUT("/:id", auth.AuthMiddleware(authService), projectHandler.UpdateProject)
		// projectRoutes.DELETE("/:id", auth.AuthMiddleware(authService), projectHandler.DeleteProject)

		// stats
		userRoutes.GET("/:id/stats", auth.AuthMiddleware(authService), statsHandler.GetUserStats)

		//team routes
		teamRoutes.POST("/", auth.AuthMiddleware(authService), teamHandler.CreateTeam)
		teamRoutes.GET("/", auth.AuthMiddleware(authService), teamHandler.GetTeamsByUserID)
		teamRoutes.DELETE("/:team_id", auth.AuthMiddleware(authService), teamHandler.DeleteTeam)

		//dashboard routes
		// get user owned teams along with team members and project assigned to each team
		// dashboardRoutes.GET("/teams", auth.AuthMiddleware(authService), dashboardService.GetUserTeams)
		// contributors routes
		contributorRoutes.POST("/", auth.AuthMiddleware(authService), contributorHandler.CreateContributor)
		contributorRoutes.GET("/", auth.AuthMiddleware(authService), contributorHandler.GetContributorsByUserID)
		contributorRoutes.GET("/:id", auth.AuthMiddleware(authService), contributorHandler.GetContributorByID)
	}

	// Start the server
	r.Run()

}
