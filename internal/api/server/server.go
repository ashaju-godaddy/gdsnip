package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ashaju-godaddy/gdsnip/internal/api/config"
	"github.com/ashaju-godaddy/gdsnip/internal/api/handlers"
	"github.com/ashaju-godaddy/gdsnip/internal/api/middleware"
	"github.com/ashaju-godaddy/gdsnip/internal/api/repository"
	"github.com/ashaju-godaddy/gdsnip/internal/api/service"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// Server represents the API server
type Server struct {
	echo   *echo.Echo
	config *config.Config
	db     *sqlx.DB
}

// New creates a new API server
func New(cfg *config.Config) (*Server, error) {
	// Initialize database
	db, err := repository.NewDB(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepo(db)
	snippetRepo := repository.NewSnippetRepo(db)
	teamRepo := repository.NewTeamRepo(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
	teamService := service.NewTeamService(teamRepo, userRepo)
	snippetService := service.NewSnippetService(snippetRepo, userRepo, teamRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	snippetHandler := handlers.NewSnippetHandler(snippetService)
	teamHandler := handlers.NewTeamHandler(teamService, snippetService)

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Setup middleware
	e.Use(echomiddleware.RequestLogger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Rate limiting
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit)
	e.Use(rateLimiter.Middleware())

	// JWT middleware for protected routes
	jwtMiddleware := middleware.JWTMiddleware(authService)
	optionalJwtMiddleware := middleware.OptionalJWTMiddleware(authService)

	// Register routes
	setupRoutes(e, authHandler, snippetHandler, teamHandler, jwtMiddleware, optionalJwtMiddleware)

	return &Server{
		echo:   e,
		config: cfg,
		db:     db,
	}, nil
}

// setupRoutes registers all API routes
func setupRoutes(e *echo.Echo, authHandler *handlers.AuthHandler, snippetHandler *handlers.SnippetHandler, teamHandler *handlers.TeamHandler, jwtMiddleware echo.MiddlewareFunc, optionalJwtMiddleware echo.MiddlewareFunc) {
	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "healthy",
		})
	})

	// API v1 group
	v1 := e.Group("/v1")

	// Public auth routes
	auth := v1.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.GET("/me", authHandler.Me, jwtMiddleware)

	// Public snippet routes
	snippets := v1.Group("/snippets")
	snippets.GET("", snippetHandler.Search)                                             // Search public snippets
	snippets.GET("/:namespace/:slug", snippetHandler.Get)                               // Get snippet details
	snippets.POST("/:namespace/:slug/pull", snippetHandler.Pull, optionalJwtMiddleware) // Pull snippet (render) - optional auth

	// Protected snippet routes
	snippets.POST("", snippetHandler.Create, jwtMiddleware) // Create snippet

	// User routes (protected)
	users := v1.Group("/users", jwtMiddleware)
	users.GET("/me/snippets", snippetHandler.ListMine) // List my snippets

	// Team routes (protected)
	teams := v1.Group("/teams", jwtMiddleware)
	teams.POST("", teamHandler.Create)                                // Create team
	teams.GET("", teamHandler.List)                                   // List my teams
	teams.GET("/:slug", teamHandler.Get)                              // Get team
	teams.PUT("/:slug", teamHandler.Update)                           // Update team
	teams.DELETE("/:slug", teamHandler.Delete)                        // Delete team
	teams.GET("/:slug/members", teamHandler.ListMembers)              // List members
	teams.POST("/:slug/members", teamHandler.AddMember)               // Add member
	teams.PUT("/:slug/members/:username", teamHandler.UpdateMemberRole) // Update role
	teams.DELETE("/:slug/members/:username", teamHandler.RemoveMember)  // Remove member
	teams.POST("/:slug/leave", teamHandler.Leave)                     // Leave team
	teams.GET("/:slug/snippets", teamHandler.ListSnippets)            // List team snippets
}

// Start starts the API server
func (s *Server) Start() error {
	addr := "http://localhost:" + s.config.Port
	fmt.Printf("API server starting on %s\n", addr)
	return s.echo.Start(":" + s.config.Port)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("🛑 Shutting down server...")

	// Close database connection
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			fmt.Printf("Error closing database: %v\n", err)
		}
	}

	// Shutdown Echo server
	return s.echo.Shutdown(ctx)
}
