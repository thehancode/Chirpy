package handlers

import (
	"net/http"
)

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
