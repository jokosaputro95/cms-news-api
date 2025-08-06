package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/jokosaputro95/cms-news-api/configs"
	"github.com/jokosaputro95/cms-news-api/internal/modules/auth/interface/rest/handlers"
)

// SetupRoutes configures all application routes
func SetupRoutes(mux *http.ServeMux, config *configs.Configs, db *sql.DB, authHandler *handlers.AuthHandler) {
	// Setup Auth routes
	SetupAuthRoutes(mux, authHandler)

	// Setup Health routes
	SetupHealthRoutes(mux, config, db)

	log.Println("✅ Routes configured successfully")
}