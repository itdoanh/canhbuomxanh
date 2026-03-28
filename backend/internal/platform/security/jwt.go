package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint64
	Role   string
}

func GenerateJWT(userID uint64, role, secret string, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}
	return signed, nil
}

func ParseJWT(tokenString, secret string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("invalid signing algorithm")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	userIDFloat, ok := claimsMap["sub"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing sub claim")
	}
	role, ok := claimsMap["role"].(string)
	if !ok {
		return nil, fmt.Errorf("missing role claim")
	}

	return &Claims{
		UserID: uint64(userIDFloat),
		Role:   role,
	}, nil
}
