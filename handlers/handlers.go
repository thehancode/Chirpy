package handlers

import (
	"Chirpy/handlers/models"
	"Chirpy/internal/database"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	respondWithOk(w, http.StatusOK)
}

func (cfg *ApiConfig) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	var userBody models.UserPostRequest
	err := json.NewDecoder(r.Body).Decode(&userBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), userBody.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}
	userResponse := models.CreateUserResponse{ID: user.ID.String(), CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
	respondWithJSON(w, 201, userResponse)
}

func (cfg *ApiConfig) PostChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirpBody models.ChirpPostRequest
	err := json.NewDecoder(r.Body).Decode(&chirpBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong with the decoding ")
		return
	}
	userID, err := uuid.Parse(chirpBody.UserId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Uuid cannot be parsed")
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

func (cfg *ApiConfig) DeleteUsersHandler(w http.ResponseWriter, r *http.Request) {
	err := cfg.db.DeleteAllUsers(r.Context())
	if cfg.platform != "dev" {
		respondWithCode(w, 403)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "something went wrong")
		return
	}
	respondWithOk(w, 200)
}
