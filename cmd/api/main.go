package main

import (
	"log"
	"oat431/fluffy-mouton/internal/bootstrap"
	"oat431/fluffy-mouton/internal/config"
	"oat431/fluffy-mouton/internal/router"
	"oat431/fluffy-mouton/pkg/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
)

func main() {
	config.LoadEnvConfig()

	// Initialize Keycloak JWKS key set (fetches public keys on startup)
	if err := utils.InitJWKS(); err != nil {
		log.Printf("WARNING: Failed to initialize JWKS: %v", err)
		log.Println("JWT validation via Authorization header will fail. Gateway header mode will still work.")
	}

	db := config.StartDatabase()
	defer db.Close()

	apiContainer, err := bootstrap.NewAPIContainer(db)
	if err != nil {
		log.Fatalf("failed to initialize api container: %v", err)
	}

	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	router.SetupRoutes(app, apiContainer, db)

	// Graceful shutdown: listen for SIGINT/SIGTERM, then drain in-flight requests
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		port := os.Getenv("PORT")
		log.Printf("starting server on :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down server gracefully...")

	if err := app.ShutdownWithTimeout(30 * time.Second); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	log.Println("server stopped cleanly")
}
