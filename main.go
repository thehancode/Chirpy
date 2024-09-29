package main

import (
	"Chirpy/handlers"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()
	apiCfg := handlers.NewApiConfig()
	fileServer := http.FileServer(http.Dir("."))
	fileserverHandler := http.StripPrefix("/app", fileServer)
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(fileserverHandler))
	mux.HandleFunc("GET /api/healthz", handlers.HealthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /api/users", apiCfg.PostUserHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.DeleteUsersHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.PostChirpHandler)

	mux.HandleFunc("GET /api/reset", apiCfg.ResetHandler)
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	log.Println("Server starting on localhost:8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
