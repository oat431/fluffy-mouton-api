package middleware

import (
	"context"
	"oat431/fluffy-mouton/pkg/common"
	"oat431/fluffy-mouton/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// OAuthMiddleware validates requests through either:
//  1. Gateway headers (X-User-Id) — when behind Flowero Gate
//  2. Direct Keycloak JWT — for direct/internal service access
//
// Sets c.Locals("user_id") (uuid.UUID), c.Locals("user_name") (string),
// and c.Locals("user_roles") ([]string).
func OAuthMiddleware(c fiber.Ctx) error {
	// ---- Path 1: Behind Gateway (trust X-User-Id header) ----
	userIDHeader := c.Get("X-User-Id")
	if userIDHeader != "" {
		userID, err := uuid.Parse(userIDHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(common.ResponseDTO[any]{
				Status: common.ERROR,
				Error: &common.ResponseDTOError{
					HttpCode:  fiber.StatusUnauthorized,
					ErrorCode: "UNAUTHORIZED",
					Message:   "Invalid X-User-Id header format",
				},
			})
		}

		c.Locals("user_id", userID)
		c.Locals("user_name", c.Get("X-User-Name", ""))

		rolesHeader := c.Get("X-User-Roles", "")
		var roles []string
		if rolesHeader != "" {
			roles = strings.Split(rolesHeader, ",")
		}
		c.Locals("user_roles", roles)

		return c.Next()
	}

	// ---- Path 2: Direct JWT validation ----
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusUnauthorized,
				ErrorCode: "UNAUTHORIZED",
				Message:   "Missing Authorization header or X-User-Id",
			},
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return c.Status(fiber.StatusUnauthorized).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusUnauthorized,
				ErrorCode: "UNAUTHORIZED",
				Message:   "Invalid Authorization header format",
			},
		})
	}

	claims, err := utils.ValidateKeycloakJWT(context.Background(), parts[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusUnauthorized,
				ErrorCode: "UNAUTHORIZED",
				Message:   "Invalid or expired token: " + err.Error(),
			},
		})
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusUnauthorized,
				ErrorCode: "UNAUTHORIZED",
				Message:   "Invalid user ID in token",
			},
		})
	}

	c.Locals("user_id", userID)
	c.Locals("user_name", claims.PreferredUsername)
	c.Locals("user_roles", claims.RealmAccess.Roles)

	return c.Next()
}

// RequireRole returns a middleware that checks for a specific role.
func RequireRole(role string) fiber.Handler {
	return func(c fiber.Ctx) error {
		roles, ok := c.Locals("user_roles").([]string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(common.ResponseDTO[any]{
				Status: common.ERROR,
				Error: &common.ResponseDTOError{
					HttpCode:  fiber.StatusForbidden,
					ErrorCode: "FORBIDDEN",
					Message:   "Access denied",
				},
			})
		}
		for _, r := range roles {
			if r == role {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusForbidden,
				ErrorCode: "FORBIDDEN",
				Message:   "Insufficient permissions",
			},
		})
	}
}

// OptionalAuth attempts to extract user info but does NOT reject
// unauthenticated requests. Useful for public endpoints.
func OptionalAuth(c fiber.Ctx) error {
	userIDHeader := c.Get("X-User-Id")
	if userIDHeader != "" {
		if userID, err := uuid.Parse(userIDHeader); err == nil {
			c.Locals("user_id", userID)
		}
	}
	return c.Next()
}
