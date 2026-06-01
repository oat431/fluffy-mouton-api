package router

import (
	"oat431/fluffy-mouton/pkg/common"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
)

type healthStatus struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Database  string `json:"database"`
}

// RegisterHealthRoutes registers a deep health check endpoint that verifies
// the database connection, not just the HTTP server is alive.
func RegisterHealthRoutes(router fiber.Router, db *sqlx.DB) {
	route := router.Group("/health")

	route.Get("/check", func(c fiber.Ctx) error {
		dbStatus := "healthy"
		if err := db.Ping(); err != nil {
			dbStatus = "unhealthy"
		}

		overallStatus := "ok"
		httpCode := fiber.StatusOK
		if dbStatus != "healthy" {
			overallStatus = "degraded"
			httpCode = fiber.StatusServiceUnavailable
		}

		return c.Status(httpCode).JSON(common.ResponseDTO[healthStatus]{
			Status: common.SUCCESS,
			Data: &healthStatus{
				Status:    overallStatus,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Database:  dbStatus,
			},
		})
	})
}
