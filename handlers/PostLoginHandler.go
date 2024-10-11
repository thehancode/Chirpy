package handlers

import (
	"Chirpy/handlers/models"
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"encoding/json"
	"net/http"
	"time"
)

func (cfg *ApiConfig) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginPostRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	loginReq.SetDefaults()

	user, err := cfg.db.GetUserByEmail(r.Context(), loginReq.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(loginReq.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Duration(*loginReq.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{Token: refreshToken, UserID: user.ID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Prepare response
	userResponse := models.CreateUserResponse{
		ID:           user.ID.String(),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	}
	respondWithJSON(w, http.StatusOK, userResponse)
}
