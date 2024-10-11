package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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
