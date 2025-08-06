package routes

import (
	"net/http"

	handlers "github.com/jokosaputro95/cms-news-api/internal/modules/auth/interface/rest/handlers"
)

func SetupAuthRoutes(mux *http.ServeMux, authHandler *handlers.AuthHandler) {
	// Auth endpoints
	mux.HandleFunc("/api/v1/auth/register", authHandler.Register)
	
	// Future auth endpoints
	// mux.HandleFunc("/api/auth/login", authHandler.Login)
	// mux.HandleFunc("/api/auth/refresh", authHandler.RefreshToken)
	// mux.HandleFunc("/api/auth/logout", authHandler.Logout)
}