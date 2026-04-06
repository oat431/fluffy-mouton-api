package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	AuthID uuid.UUID `json:"auth_id"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("SECRET"))

func GenerateAccessToken(authID uuid.UUID) (string, error) {
	claims := JWTClaims{
		AuthID: authID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // 1 hour
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "fluffy-mouton",
			Subject:   authID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken() (string, error) {
	// Simple random string for refresh token, could also be a JWT
	return GenerateUUID(), nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
