package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) PostPolkaWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := cfg.AuthenticateApiKey(r)
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid authorization token")
		return
	}
	var webhookReq struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	err = json.NewDecoder(r.Body).Decode(&webhookReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if webhookReq.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Parse the user ID
	userID, err := uuid.Parse(webhookReq.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		log.Printf("Error updating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to upgrade user")
		return
	}

	// Successfully upgraded user, respond with 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
