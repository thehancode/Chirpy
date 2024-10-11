package handlers

import (
	"Chirpy/handlers/models"
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (cfg *ApiConfig) PutUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.AuthenticateRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	var userBody models.UserPutRequest
	err = json.NewDecoder(r.Body).Decode(&userBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong decoding")
		return
	}
	hashedPasswordRequest, err := auth.HashPassword(userBody.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Cannot hash the password ")
		return
	}
	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{ID: userID, Email: userBody.Email, HashedPassword: hashedPasswordRequest})
	fmt.Printf("Updated user: %+v\n", user)
	if err != nil {
		log.Printf("Error creating user: %v", err) // Log the SQL error
		respondWithError(w, http.StatusBadRequest, "Something went wrong creating")
		return
	}

	userResponse := models.CreateUserResponse{ID: user.ID.String(), CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, IsChirpyRed: user.IsChirpyRed}
	respondWithJSON(w, 200, userResponse)
}
