package main

import (
	"fmt"
	"net/http"
	"os"
)

type ApiConfig struct {
	fileserverHits int
}

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits = 0
	hits := fmt.Sprintf("Hits %d", cfg.fileserverHits)
	w.Write([]byte(hits))
}

func (cfg *ApiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
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
