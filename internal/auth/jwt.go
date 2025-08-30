package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	key := []byte(tokenSecret)
	issuedAt := jwt.NewNumericDate(time.Now())
	expiresAt := jwt.NewNumericDate(time.Now().Add(expiresIn))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		Subject:   userID.String(),
	})

	return token.SignedString(key)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("Could not parse with claims: %v", err)
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, fmt.Errorf("Couldn't get the user ID from claims: %v", err)
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, fmt.Errorf("Couldn't get the issuer: %v", err)
	}
	if issuer != "chirpy" {
		return uuid.Nil, fmt.Errorf("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Couldn't parse the user ID: %v", err)
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	var token string
	authHeaderParts := strings.SplitN(headers.Get("Authorization"), " ", 2)
	if len(authHeaderParts) == 2 && authHeaderParts[0] == "Bearer" {
		token = authHeaderParts[1]
	} else {
		return "", fmt.Errorf("Couldn't find token string")
	}

	return token, nil
}
