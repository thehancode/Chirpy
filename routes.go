package main

import (
	"Chirpy/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter sets up the HTTP routes using Gorilla Mux and returns the router.
func NewRouter(apiCfg *handlers.ApiConfig) http.Handler {
	router := mux.NewRouter()

	// Serve static files under /app/
	fileServer := http.FileServer(http.Dir("."))
	fileserverHandler := http.StripPrefix("/app", fileServer)
	router.PathPrefix("/app/").Handler(apiCfg.MiddlewareMetricsInc(fileserverHandler))

	// API routes
	router.HandleFunc("/api/healthz", handlers.HealthzHandler).Methods("GET")
	router.HandleFunc("/admin/metrics", apiCfg.MetricsHandler).Methods("GET")
	router.HandleFunc("/api/login", apiCfg.PostLoginHandler).Methods("POST")
	router.HandleFunc("/api/refresh", apiCfg.PostRefreshHandler).Methods("POST")
	router.HandleFunc("/api/revoke", apiCfg.PostRevokeHandler).Methods("POST")
	router.HandleFunc("/api/users", apiCfg.PostUserHandler).Methods("POST")
	router.HandleFunc("/api/users", apiCfg.GetUsersHandler).Methods("GET")
	router.HandleFunc("/api/users", apiCfg.PutUserHandler).Methods("PUT")
	router.HandleFunc("/admin/reset", apiCfg.DeleteUsersHandler).Methods("POST")
	router.HandleFunc("/api/chirps", apiCfg.PostChirpHandler).Methods("POST")
	router.HandleFunc("/api/chirps", apiCfg.GetChirpsHandler).Methods("GET")
	router.HandleFunc("/api/chirps/{id}", apiCfg.GetChirpHandler).Methods("GET")

	return router
}
