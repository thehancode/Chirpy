package handlers

import (
	"Chirpy/handlers/models"
	"Chirpy/internal/database"
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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
