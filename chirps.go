package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/bencuci/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerPostChirp(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	// in case we cannot decode the response
	if err := decoder.Decode(&params); err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Could not decode request", err)
		return
	}

	// in case response body length exceeds the limit
	err := handlerValidateChirp(rw, params.Body)
	if err != nil {
		respondWithError(rw, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	createdChirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   getCleanedBody(params.Body),
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Could not post the chirp", err)
		return
	}

	chirp := Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}

	respondWithJSON(rw, http.StatusCreated, chirp)
}

func handlerValidateChirp(rw http.ResponseWriter, chirpBody string) error {
	const maxChirpLength = 140
	if len(chirpBody) > maxChirpLength {
		return errors.New("Chirp is too long")
	}

	return nil
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
