package handlers

import (
	"Chirpy/handlers/models"
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *ApiConfig) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	var userBody models.UserPostRequest
	err := json.NewDecoder(r.Body).Decode(&userBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong decoding")
		return
	}
	hashedPasswordRequest, err := auth.HashPassword(userBody.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Cannot hash the password ")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: userBody.Email, HashedPassword: hashedPasswordRequest})
	if err != nil {
		log.Printf("Error creating user: %v", err) // Log the SQL error
		respondWithError(w, http.StatusBadRequest, "Something went wrong creating")
		return
	}
	userResponse := models.CreateUserResponse{ID: user.ID.String(), CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, IsChirpyRed: user.IsChirpyRed}
	respondWithJSON(w, 201, userResponse)
}
