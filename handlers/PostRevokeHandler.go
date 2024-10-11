package handlers

import (
	"Chirpy/internal/auth"
	"log"
	"net/http"
	"time"
)

func (cfg *ApiConfig) PostRevokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid authorization token")
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
		respondWithError(w, http.StatusUnauthorized, "Refresh token has already been revoked")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
