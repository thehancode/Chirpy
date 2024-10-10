package handlers

import (
	"Chirpy/internal/database"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ApiConfig struct {
	fileserverHits int
	db             *database.Queries
	platform       string
	tokenSecret    string
	polkaKey       string
}

func NewApiConfig() *ApiConfig {
	err := godotenv.Load()
	if err != nil {
		return nil
	}
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECTET")

	db, err := sql.Open("postgres", dbURL)
	polkaKey := os.Getenv("POLKA_KEY")
	if err != nil {
		return nil
	}
	dbQueries := database.New(db)

	return &ApiConfig{
		fileserverHits: 0,
		db:             dbQueries,
		platform:       platform,
		tokenSecret:    secret,
		polkaKey:       polkaKey,
	}
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits = 0
	hits := fmt.Sprintf("Hits %d", cfg.fileserverHits)
	w.Write([]byte(hits))
}

func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	templateContent, err := os.ReadFile("metrics.html")
	if err != nil {
		http.Error(w, "Error reading template file", http.StatusInternalServerError)
		return
	}
	templateString := string(templateContent)
	htmlContent := fmt.Sprintf(templateString, cfg.fileserverHits)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlContent))
}
