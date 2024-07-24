package main

import (
	"Chirpy/handlers/models"
	"encoding/json"
	"net/http"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	respondWithOk(w, http.StatusOK)
}

func (cfg *ApiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirp models.ChirpRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	respondWithJSON(w, http.StatusOK, models.ValidResponse{Valid: true})
}
