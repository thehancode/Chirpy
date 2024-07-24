package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	apiCfg := ApiConfig{fileserverHits: 0}
	fileServer := http.FileServer(http.Dir("."))
	fileserverHandler := http.StripPrefix("/app", fileServer)
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(fileserverHandler))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirpHandler)

	mux.HandleFunc("GET /api/reset", apiCfg.resetHandler)
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
