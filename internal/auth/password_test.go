package auth

import "testing"

func TestCheckPasswordHash(t *testing.T) {
	password := "password1234567"
	hash, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}

	if err := CheckPasswordHash(password, hash); err != nil {
		t.Error(err)
	}
}
