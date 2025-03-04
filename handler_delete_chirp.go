package main

import (
	"net/http"

	"github.com/Guaty/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid user or token", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "you may not delete this chirp", err)
		return
	}

	if err = cfg.db.DeleteChirp(r.Context(), chirpID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
