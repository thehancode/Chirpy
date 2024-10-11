package handlers

import (
	"Chirpy/handlers/models"
	"Chirpy/internal/database"
	"encoding/json"
	"net/http"
)

func (cfg *ApiConfig) PostChirpHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.AuthenticateRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	var chirpBody models.ChirpPostRequest
	err = json.NewDecoder(r.Body).Decode(&chirpBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong with the decoding ")
		return
	}
	if len(chirpBody.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	validated := validateChirp(chirpBody.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{Body: validated, UserID: userID})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong inserting chirp "+err.Error())
		return
	}
	chirpResponse := models.CreateChirpResponse{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	}
	respondWithJSON(w, 201, chirpResponse)
}
