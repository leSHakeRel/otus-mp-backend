package main

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"movie-night-planner-backend/internal/config"
	"movie-night-planner-backend/internal/database"
	"movie-night-planner-backend/internal/handlers"
	"movie-night-planner-backend/internal/middleware"
	"movie-night-planner-backend/internal/repositories"
	"movie-night-planner-backend/internal/services"
	"movie-night-planner-backend/internal/tmdb"
	"movie-night-planner-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	err = database.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.DB)
	eveningRepo := repositories.NewEveningRepository(database.DB)
	eveningFilmRepo := repositories.NewEveningFilmRepository(database.DB)
	voteRepo := repositories.NewVoteRepository(database.DB)
	commentRepo := repositories.NewCommentRepository(database.DB)

	// Initialize services
	jwtService := utils.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)
	authService := services.NewAuthService(userRepo, jwtService)
	eveningService := services.NewEveningService(eveningRepo, eveningFilmRepo, voteRepo, commentRepo, userRepo)
	tmdbClient := tmdb.NewClient(&cfg.TMDB)
	movieService := services.NewMovieService(eveningFilmRepo, eveningRepo, tmdbClient)
	voteService := services.NewVoteService(voteRepo, eveningRepo, eveningFilmRepo, userRepo)
	commentService := services.NewCommentService(commentRepo, eveningRepo, userRepo)

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	eveningHandler := handlers.NewEveningHandler(eveningService)
	movieHandler := handlers.NewMovieHandler(movieService)
	voteHandler := handlers.NewVoteHandler(voteService)
	commentHandler := handlers.NewCommentHandler(commentService)
	userHandler := handlers.NewUserHandler(userService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	corsMiddleware := middleware.CORSMiddleware(cfg.CORS.AllowedOrigins)

	// Setup Gin router
	r := gin.Default()
	r.Use(corsMiddleware)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes
	api := r.Group("/api/v1")
	{
		// Auth routes (no auth required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(authMiddleware.Authenticate())
		{
			protected.GET("/auth/me", authHandler.GetMe)
		}

		// Evening routes (public read only, optional auth for filter=my)
		evenings := api.Group("/evenings")
		evenings.Use(authMiddleware.OptionalAuth())
		{
			evenings.GET("", eveningHandler.GetAllEvenings)
			evenings.GET("/:id", eveningHandler.GetEvening)
			evenings.GET("/:id/movies", movieHandler.GetFilmsForEvening)
			evenings.GET("/:id/votes", voteHandler.GetVotesForEvening)
			evenings.GET("/:id/comments", commentHandler.GetCommentsForEvening)
		}

		// Evening routes (protected - write operations)
		protectedEvenings := api.Group("/evenings")
		protectedEvenings.Use(authMiddleware.Authenticate())
		{
			protectedEvenings.POST("", eveningHandler.CreateEvening)
			protectedEvenings.PUT("/:id", eveningHandler.UpdateEvening)
			protectedEvenings.DELETE("/:id", eveningHandler.DeleteEvening)

			// Evening films
			protectedEvenings.POST("/:id/movies", movieHandler.AddFilmToEvening)
			protectedEvenings.DELETE("/:id/movies/:tmdbId", movieHandler.RemoveFilmFromEvening)

			// Evening votes
			protectedEvenings.POST("/:id/votes", voteHandler.CreateVote)
			protectedEvenings.DELETE("/:id/votes/:voteId", voteHandler.DeleteVote)

			// Evening comments
			protectedEvenings.POST("/:id/comments", commentHandler.CreateComment)
		}

		// User routes (protected)
		users := api.Group("/users")
		users.Use(authMiddleware.Authenticate())
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:userId", userHandler.GetUser)
		}

		// Movie routes (public)
		movies := api.Group("/movies")
		{
			movies.GET("/popular", movieHandler.GetPopularMovies)
			movies.GET("/search", movieHandler.SearchMovies)
			movies.GET("/:tmdbId", movieHandler.GetMovieDetails)
		}
	}

	// Start server
	addr := ":" + cfg.Server.Port
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
