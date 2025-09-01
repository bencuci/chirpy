package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
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
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirps(rw http.ResponseWriter, req *http.Request) {
	sortMethod := req.URL.Query().Get("sort")

	var err error
	authorID := uuid.Nil
	authorIDString := req.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(rw, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
	}

	chirps := []database.Chirp{}
	if authorID != uuid.Nil {
		chirps, err = cfg.dbQueries.GetChirpsFromUserID(req.Context(), authorID)
		if err != nil {
			respondWithError(rw, http.StatusInternalServerError, "Could not get chirps", err)
			return
		}
	} else {
		chirps, err = cfg.dbQueries.GetChirps(req.Context())
		if err != nil {
			respondWithError(rw, http.StatusInternalServerError, "Could not get chirps", err)
			return
		}
	}

	chirpsResponse := []Chirp{}
	for _, c := range chirps {
		chirp := Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
		chirpsResponse = append(chirpsResponse, chirp)
	}

	if sortMethod == "desc" {
		sort.Slice(chirpsResponse, func(i, j int) bool {
			return chirpsResponse[i].CreatedAt.After(chirpsResponse[j].CreatedAt)
		})
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
		respondWithError(rw, http.StatusNotFound, "Couldn't find chirp", err)
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

	// in case we cannot decode the response
	if err := decoder.Decode(&params); err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Could not decode request", err)
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, err.Error(), err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, err.Error(), err)
		return
	}

	// in case response body length exceeds the limit
	err = validateChirp(params.Body)
	if err != nil {
		respondWithError(rw, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	createdChirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   getCleanedBody(params.Body),
		UserID: userID,
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
		UserID:    userID,
	}

	respondWithJSON(rw, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerDeleteChirp(rw http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, err.Error(), err)
		return
	}

	tokenUserID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(rw, http.StatusUnauthorized, err.Error(), err)
		return
	}

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't parse path value", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), chirpID)
	if err != nil {
		respondWithError(rw, http.StatusNotFound, "Couldn't find chirp", err)
		return
	}

	if chirp.UserID != tokenUserID {
		respondWithError(rw, http.StatusForbidden, "Operation forbidden", nil)
		return
	}

	cfg.dbQueries.DeleteChirp(req.Context(), chirp.ID)

	respondWithJSON(rw, http.StatusNoContent, nil)
}

func validateChirp(chirpBody string) error {
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
