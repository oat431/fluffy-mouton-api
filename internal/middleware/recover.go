package middleware

import (
	"oat431/fluffy-mouton/pkg/common"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

// Recover is a panic recovery middleware that catches panics in handlers
// and returns a 500 Internal Server Error instead of crashing the server.
func Recover(c fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("panic recovered: %v", r)
			c.Status(fiber.StatusInternalServerError).JSON(common.ResponseDTO[any]{
				Status: common.ERROR,
				Error: &common.ResponseDTOError{
					HttpCode:  fiber.StatusInternalServerError,
					ErrorCode: "INTERNAL_SERVER_ERROR",
					Message:   "An unexpected error occurred",
				},
			})
		}
	}()

	return c.Next()
}
