package handlers

import (
	"Chirpy/handlers/models"
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	respondWithOk(w, http.StatusOK)
}

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

func (cfg *ApiConfig) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var userList []database.User
	userList, err := cfg.db.GetAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong inserting user "+err.Error())
		return
	}
	userResponses := make([]models.CreateUserResponse, len(userList))

	for i, user := range userList {
		userResponses[i] = models.CreateUserResponse{ID: user.ID.String(), CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
	}
	respondWithJSON(w, 200, userResponses)
}

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

func (cfg *ApiConfig) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the author_id from query parameters
	authorIDStr := r.URL.Query().Get("author_id")
	sortDirection := r.URL.Query().Get("sort")

	if sortDirection != "asc" && sortDirection != "desc" {
		sortDirection = "asc"
	}

	var chirpList []database.Chirp
	var err error

	if authorIDStr == "" {
		chirpList, err = cfg.db.GetAllChirps(r.Context(), sortDirection)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error fetching chirps: "+err.Error())
			return
		}
	} else {
		authorID, err := uuid.Parse(authorIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id: "+err.Error())
			return
		}
		chirpList, err = cfg.db.GetChirpsByAuthorID(r.Context(), authorID, sortDirection)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error fetching chirps for author: "+err.Error())
			return
		}
	}

	// Transform chirpList into response format
	chirpResponses := make([]models.CreateChirpResponse, len(chirpList))
	for i, chirp := range chirpList {
		chirpResponses[i] = models.CreateChirpResponse{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID.String(),
		}
	}

	// Respond with the list of chirps
	respondWithJSON(w, http.StatusOK, chirpResponses)
}

func (cfg *ApiConfig) GetChirpHandler(w http.ResponseWriter, r *http.Request) {
	var chirp database.Chirp
	vars := mux.Vars(r)
	id := vars["id"]
	chirpId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid UUID: id cannot be parsed as uuid")
		return
	}
	chirp, err = cfg.db.GetChirp(r.Context(), chirpId)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	if err != nil {
		log.Printf("Error retrieving chirp with ID %s: %v", chirpId, err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong while inserting chirp "+err.Error())
		return
	}
	chirpResponse := models.CreateChirpResponse{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID.String(),
	}
	respondWithJSON(w, 200, chirpResponse)
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

func (cfg *ApiConfig) DeleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.AuthenticateRequest(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	chirpId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID: id cannot be parsed as UUID")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	if err != nil {
		log.Printf("Error retrieving chirp with ID %s: %v", chirpId, err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while retrieving chirp: "+err.Error())
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp")
		return
	}

	_, err = cfg.db.DeleteChirp(r.Context(), chirpId)
	if err != nil {
		log.Printf("Error deleting chirp with ID %s: %v", chirpId, err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong while deleting chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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
