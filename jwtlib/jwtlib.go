package jwtlib

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTHandler struct {
	secret []byte
}

// New creates a new instance of the JWT handler with a secret key.
func New(secret string) *JWTHandler {
	return &JWTHandler{secret: []byte(secret)}
}

// Generate creates a JWT with custom claims.
func (j *JWTHandler) Generate(claims map[string]interface{}, expiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(expiry).Unix(),
		"data": claims,
	})
	return token.SignedString(j.secret)
}

// Verify checks if a token is valid and extracts claims.
func (j *JWTHandler) Verify(tokenString string) (bool, map[string]interface{}) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil || !token.Valid {
		return false, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	return true, claims["data"].(map[string]interface{})
}
