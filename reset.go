package main

import (
	"errors"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "you don't have permission for this", errors.New("unauthorized"))
		return
	}

	if err := cfg.db.DeleteUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete users", err)
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0, database reset to initial state"))
}
