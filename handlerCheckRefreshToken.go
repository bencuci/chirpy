package main

import (
	"net/http"
	"time"

	"github.com/bencuci/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(rw http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, err.Error(), err)
		return
	}

	user, err := cfg.dbQueries.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, err.Error(), err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, 1*time.Hour)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Could not create jwt token: %v", err)
		return
	}

	respondWithJSON(rw, http.StatusOK, response{Token: accessToken})
}

func (cfg *apiConfig) handlerRevokeRefreshToken(rw http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, err.Error(), err)
		return
	}

	_, err = cfg.dbQueries.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, "Couldn't find the refresh token", err)
	}

	cfg.dbQueries.RevokeRefreshToken(req.Context(), refreshToken)

	respondWithJSON(rw, http.StatusNoContent, nil)
}
