package handlers

import (
	"Chirpy/internal/auth"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (cfg *ApiConfig) PostRefreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Errorf("missing or invalid authorization token")
		return
	}
	refToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Error retrieving refresh token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	if time.Now().After(refToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has expired")
		return
	}

	if refToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has been revoked")
		return
	}

	newAccessToken, err := auth.MakeJWT(refToken.UserID, cfg.tokenSecret, time.Hour)
	if err != nil {
		log.Printf("Error generating new access token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	response := map[string]string{
		"token": newAccessToken,
	}
	respondWithJSON(w, http.StatusOK, response)
}
