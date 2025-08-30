package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bencuci/chirpy/internal/auth"
	"github.com/bencuci/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handlerGetChirps(rw http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Could not get chirps", err)
		return
	}

	chirpsResponse := []Chirp{}
	for _, chirpFromDB := range chirps {
		chirp := Chirp{
			ID:        chirpFromDB.ID,
			CreatedAt: chirpFromDB.CreatedAt,
			UpdatedAt: chirpFromDB.UpdatedAt,
			Body:      chirpFromDB.Body,
			UserID:    chirpFromDB.UserID,
		}
		chirpsResponse = append(chirpsResponse, chirp)
	}

	respondWithJSON(rw, http.StatusOK, chirpsResponse)
}

func (cfg *apiConfig) handlerGetChirp(rw http.ResponseWriter, req *http.Request) {
	userID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't parse path value", err)
		return
	}
	chirp, err := cfg.dbQueries.GetChirp(req.Context(), userID)
	if err != nil {
		respondWithError(rw, http.StatusNotFound, "Couldn't find user", err)
		return
	}

	respondWithJSON(rw, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerPostChirp(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, http.StatusBadRequest, err.Error(), err)
	}
	_, err = auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, err.Error(), err)
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(rw, http.StatusBadRequest, err.Error(), err)
		return
	}

	fmt.Printf("Requesting with ID of: %s", params.UserID)
	chirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(rw, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
