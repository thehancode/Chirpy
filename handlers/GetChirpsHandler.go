package handlers

import (
	"Chirpy/handlers/models"
	"Chirpy/internal/database"
	"net/http"

	"github.com/google/uuid"
)

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
