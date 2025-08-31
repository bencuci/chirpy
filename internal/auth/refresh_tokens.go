package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	rand32bytes := make([]byte, 32)
	_, err := rand.Read(rand32bytes)

	return hex.EncodeToString(rand32bytes), err
}
