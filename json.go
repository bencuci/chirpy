package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(rw http.ResponseWriter, statusCode int, msg string, err error) {
	type responseError struct {
		Error string `json:"error"`
	}

	if err != nil {
		log.Println(err)
	}
	if statusCode >= 500 {
		log.Printf("Responding with 5XX status code: %s", msg)
	}

	respondWithJSON(rw, statusCode, responseError{
		Error: msg,
	})
}

func respondWithJSON(rw http.ResponseWriter, statusCode int, payload any) {
	rw.Header().Set("Content-type", "application/json")

	resp, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Something went wrong: %s", err)
		rw.WriteHeader(500)
		return
	}

	rw.WriteHeader(statusCode)
	rw.Write(resp)
}
