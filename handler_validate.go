package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidate(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type response struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	// in case we cannot decode the response
	if err := decoder.Decode(&params); err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Could not decode request", err)
		return
	}

	// in case response body length exceeds the limit
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(rw, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	resp := response{
		Valid: true,
	}
	respondWithJSON(rw, http.StatusOK, resp)
}
