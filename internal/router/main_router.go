package router

import (
	"oat431/fluffy-mouton/internal/bootstrap"
	"oat431/fluffy-mouton/internal/config"
	"oat431/fluffy-mouton/internal/middleware"

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
	app.Use(cors.New(config.InitCorsConfig()))

	api := app.Group("/api")
	v1 := api.Group("/v1")

	RegisterHealthRoutes(v1, db)
	RegisterAuthRoutes(v1, apiContainer.AuthController)
	RegisterShortLinkRoutes(v1, apiContainer.ShortLinkController)
	RegisterRedirectRoutes(v1, apiContainer.RedirectController)
}
