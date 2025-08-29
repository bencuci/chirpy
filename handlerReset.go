package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerResetVisiterCount(rw http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(rw, http.StatusForbidden, "Forbidden access", nil)
		return
	}

	err := cfg.dbQueries.ResetUsers(req.Context())
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Could not reset users", err)
		return
	}
	cfg.fileserverHits.Store(0)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Database has been reset."))
}
