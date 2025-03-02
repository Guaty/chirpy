package main

import (
	"net/http"

	"github.com/Guaty/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get token", err)
		return
	}

	err = cfg.db.RevokeToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "token not found", err)
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
