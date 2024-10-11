package handlers

import (
	"Chirpy/internal/auth"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	respondWithOk(w, http.StatusOK)
}

func (cfg *ApiConfig) AuthenticateRequest(r *http.Request) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.Nil, fmt.Errorf("missing or invalid authorization token")
	}

	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	return userID, nil
}

func (cfg *ApiConfig) AuthenticateApiKey(r *http.Request) (string, error) {
	token, err := auth.GetAPIKey(r.Header)
	if err != nil {
		return "", fmt.Errorf("missing or invalid authorization token")
	}

	return token, nil
}
