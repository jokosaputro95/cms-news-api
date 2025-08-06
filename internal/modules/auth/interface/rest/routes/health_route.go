package routes

import (
	"database/sql"
	"net/http"

	"github.com/jokosaputro95/cms-news-api/configs"
)

// SetupHealthRoutes configures health check routes
func SetupHealthRoutes(mux *http.ServeMux, config *configs.Configs, db *sql.DB) {
	// Health check endpoints
	mux.HandleFunc("/health", healthCheck(config))
	mux.HandleFunc("/health/db", dbHealthCheck(db))
}

// healthCheck returns basic health status
func healthCheck(config *configs.Configs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "ok",
			"message": "Server is running",
			"app": "` + config.AppName + `",
			"version": "` + config.AppVersion + `"
		}`))
	}
}

// dbHealthCheck returns database health status
func dbHealthCheck(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := db.Ping()
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{
				"status": "error",
				"message": "Database connection failed",
				"error": "` + err.Error() + `"
			}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "ok",
			"message": "Database is healthy"
		}`))
	}
}