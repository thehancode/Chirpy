package handlers

import (
	"Chirpy/handlers/models"
	"Chirpy/internal/database"
	"net/http"
)

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
