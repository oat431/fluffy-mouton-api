package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwksOnce sync.Once
	jwksKF   keyfunc.Keyfunc
	jwksErr  error
)

// InitJWKS fetches the Keycloak JWKS endpoint once and caches the public keys.
// Must be called once during application startup.
func InitJWKS() error {
	jwksOnce.Do(func() {
		issuerURL := os.Getenv("OAUTH_ISSUER_URI")
		if issuerURL == "" {
			jwksErr = fmt.Errorf("OAUTH_ISSUER_URI is not set")
			return
		}
		jwksURL := issuerURL + "/protocol/openid-connect/certs"

		jwksKF, jwksErr = keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
		if jwksErr != nil {
			return
		}
		log.Printf("JWKS initialized from %s", jwksURL)
	})
	return jwksErr
}

// ValidateKeycloakJWT validates a JWT against the cached Keycloak JWKS key set.
func ValidateKeycloakJWT(ctx context.Context, tokenString string) (*KeycloakClaims, error) {
	if jwksKF == nil {
		return nil, fmt.Errorf("JWKS not initialized: call InitJWKS() first")
	}

	issuerURL := os.Getenv("OAUTH_ISSUER_URI")

	token, err := jwt.ParseWithClaims(
		tokenString,
		&KeycloakClaims{},
		jwksKF.KeyfuncCtx(ctx),
		jwt.WithLeeway(5*time.Second),
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(issuerURL),
	)
	if err != nil {
		return nil, fmt.Errorf("jwt validation failed: %w", err)
	}

	claims, ok := token.Claims.(*KeycloakClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// KeycloakClaims represents the claims extracted from a Keycloak JWT.
type KeycloakClaims struct {
	jwt.RegisteredClaims

	PreferredUsername string      `json:"preferred_username"`
	Email             string      `json:"email"`
	EmailVerified     bool        `json:"email_verified"`
	RealmAccess       RealmAccess `json:"realm_access"`
}

// RealmAccess contains realm-level roles.
type RealmAccess struct {
	Roles []string `json:"roles"`
}

// HasRole checks if the JWT contains a specific realm role.
func (c *KeycloakClaims) HasRole(role string) bool {
	for _, r := range c.RealmAccess.Roles {
		if r == role {
			return true
		}
	}
	return false
}
