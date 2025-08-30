package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bencuci/chirpy/internal/auth"
	"github.com/bencuci/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	Token          string    `json:"token"`
}

func (cfg *apiConfig) handlerCreateUser(rw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Could not decode request", err)
		return
	}

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Couldn't hash the password", err)
		return
	}
	user, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPW,
	})
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, "Could not create user", err)
		return
	}

	respondWithJSON(rw, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
