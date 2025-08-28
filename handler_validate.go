package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(rw http.ResponseWriter, req *http.Request) {
	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}
	type parameters struct {
		Body string `json:"body"`
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

	// filter out banned words
	resp := response{
		CleanedBody: getCleanedBody(params.Body),
	}

	respondWithJSON(rw, http.StatusOK, resp)
}

func getCleanedBody(body string) string {
	bannedWords := map[string]struct{}{
		"kerfuffle": {}, "sharbert": {}, "fornax": {},
	}
	words := strings.Fields(body)
	for i, word := range words {
		if _, exists := bannedWords[strings.ToLower(word)]; exists {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
