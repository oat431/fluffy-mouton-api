package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

func RegisterHealthRoutes(router fiber.Router) {
	route := router.Group("/health")

	route.Get("/check", healthcheck.New())
}
