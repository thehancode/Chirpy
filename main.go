package main

import (
	"Chirpy/handlers"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	apiCfg := handlers.NewApiConfig()
	router := NewRouter(apiCfg)
	addr := ":8080"
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
