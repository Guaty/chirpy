package main

import (
	"net/http"
	"time"

	"github.com/Guaty/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get token", err)
		return
	}

	refreshToken, err := cfg.db.GetToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no vailid token found", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token for user", err)
		return
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate new token", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})
}
