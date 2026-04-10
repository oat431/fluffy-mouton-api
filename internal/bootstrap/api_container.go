package bootstrap

import (
	"oat431/fluffy-mouton/internal/config"
	"oat431/fluffy-mouton/internal/controller"
	"oat431/fluffy-mouton/internal/repository"
	"oat431/fluffy-mouton/internal/service"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jmoiron/sqlx"
)

type APIContainer struct {
	AuthController *controller.AuthController
}

func NewAPIContainer(db *sqlx.DB) *APIContainer {
	log.Info("Registering Auth Repository")
	authRepository := repository.NewAuthRepository(db)

	log.Info("Registering Refresh Token Repository")
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)

	log.Info("Registering Email Verify Token Repository")
	emailVerifyTokenRepository := repository.NewEmailVerifyTokenRepository(db)

	log.Info("Registering SMTP Service")
	smtpService := service.NewSMTPService(config.GetEmailConfig())

	log.Info("Registering Auth Service")
	authService := service.NewAuthService(authRepository, refreshTokenRepository, emailVerifyTokenRepository, smtpService)

	log.Info("Registering Auth Controller")
	authController := controller.NewAuthController(authService)

	log.Info("Registered All API")
	return &APIContainer{
		AuthController: authController,
	}
}
