package router

import (
	"oat431/fluffy-mouton/internal/bootstrap"
	"oat431/fluffy-mouton/internal/middleware"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/jmoiron/sqlx"
)

func init() {
	log.Info("Initializing routes...")
}

func SetupRoutes(app *fiber.App, apiContainer *bootstrap.APIContainer, db *sqlx.DB) {
	// Order matters: recover must be first to catch panics everywhere
	app.Use(middleware.Recover)
	app.Use(middleware.RequestID)
	app.Use(middleware.GlobalLogger)

	// CORS: enabled in local dev (behind gateway it's handled by the gateway)
	if os.Getenv("CORS_ENABLED") == "true" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Authorization", "Content-Type", "X-Requested-With"},
			AllowCredentials: true,
		}))
	}

	api := app.Group("/api")
	v1 := api.Group("/v1")

	RegisterHealthRoutes(v1, db)
	RegisterShortLinkRoutes(v1, apiContainer.ShortLinkController)
	RegisterRedirectRoutes(v1, apiContainer.RedirectController)
}
