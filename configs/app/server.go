package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/jokosaputro95/cms-news-api/configs"
	"github.com/jokosaputro95/cms-news-api/internal/modules/auth/application/usecases"
	"github.com/jokosaputro95/cms-news-api/internal/modules/auth/infrastructure/persistence/repositories"
	"github.com/jokosaputro95/cms-news-api/internal/modules/auth/infrastructure/security"
	"github.com/jokosaputro95/cms-news-api/internal/modules/auth/interface/rest/handlers"
	"github.com/jokosaputro95/cms-news-api/internal/modules/auth/interface/rest/routes"
	"github.com/jokosaputro95/cms-news-api/internal/shared"
)

type Server struct {
	config      *configs.Configs
	mux         *http.ServeMux
	db          *sql.DB
	authHandler *handlers.AuthHandler
}

func Run() {
	// Load config
	cfg := configs.LoadConfig(
		true, // false = dev, true = test
		".env",
	)

	// Create server instance
	server := &Server{
		config: cfg,
		mux:    http.NewServeMux(),
	}

	// Initialize server
	err := server.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Start HTTP server
	err = server.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) Initialize() error {
	// 1. Setup database connection
	err := s.setupDatabase()
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// 2. Setup dependencies (Dependency Injection)
	s.setupDependencies()

	// 3. Setup routes (menggunakan routes package)
	s.setupRoutes()

	log.Println("✅ Server initialized successfully")
	return nil
}

func (s *Server) setupDatabase() error {
	dsn := s.config.GetDatabaseDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// Setup connection pool
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	s.db = db
	log.Println("✅ Database connected with connection pool")
	return nil
}

func (s *Server) setupDependencies() {
	// === Infrastructure Layer ===
	// Repositories
	userRepository := repositories.NewUserRepositoryPostgres(s.db)

	// Security services  
	hasher := security.NewBcryptHasher(12)

	// Shared services
	uuidGenerator := &shared.DefaultUUIDGenerator{}

	// === Application Layer ===
	// Use Cases
	registerUserUseCase := usecases.NewRegisterUser(
		userRepository,
		uuidGenerator,
		hasher,
	)

	// === Interface Layer ===
	// Handlers
	s.authHandler = handlers.NewAuthHandler(registerUserUseCase)

	log.Println("✅ Dependencies wired successfully")
}

// ✅ Server sekarang clean - hanya delegate ke routes package
func (s *Server) setupRoutes() {
	routes.SetupRoutes(s.mux, s.config, s.db, s.authHandler)
}

func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%s", s.config.ServerHost, s.config.ServerPort)

	log.Printf("INFO: %-16s: %s", "APP_NAME", s.config.AppName)
	log.Printf("INFO: %-16s: %s", "APP_VERSION", s.config.AppVersion)
	log.Printf("INFO: %-16s: %s", "APP_ENV", s.config.AppEnv)
	log.Printf("🚀 Server starting on http://%s", address)

	return http.ListenAndServe(address, s.mux)
}