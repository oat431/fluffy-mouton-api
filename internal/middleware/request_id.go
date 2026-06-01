package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// RequestID generates a unique request ID for every incoming request
// and attaches it to both the response header and the request context.
func RequestID(c fiber.Ctx) error {
	requestID := uuid.New().String()
	c.Set("X-Request-ID", requestID)
	c.Locals("request_id", requestID)
	return c.Next()
}
