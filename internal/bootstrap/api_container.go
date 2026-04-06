package bootstrap

import (
	"github.com/gofiber/fiber/v3/log"
	"github.com/jmoiron/sqlx"
)

type APIContainer struct {
}

func NewAPIContainer(db *sqlx.DB) *APIContainer {
	log.Info("Registering Auth Repository")

	log.Info("Registering Refresh Token Repository")

	log.Info("Registering Email Verify Token Repository")

	log.Info("Registering SMTP Service")

	log.Info("Registering Auth Service")

	log.Info("Registering Auth Controller")

	log.Info("Registered All API")
	return &APIContainer{}
}
