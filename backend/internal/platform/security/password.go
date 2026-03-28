package security

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(raw string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func VerifyPassword(hashValue, raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashValue), []byte(raw)) == nil
}
