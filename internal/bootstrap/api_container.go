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
	AuthController      *controller.AuthController
	ShortLinkController *controller.ShortLinkController
	RedirectController  *controller.RedirectController
}

func NewAPIContainer(db *sqlx.DB) *APIContainer {
	log.Info("Registering Application Repository")
	authRepository := repository.NewAuthRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)
	emailVerifyTokenRepository := repository.NewEmailVerifyTokenRepository(db)
	shortLinkRepository := repository.NewShortLinkRepository(db)

	log.Info("Registering Application Service")
	smtpService := service.NewSMTPService(config.GetEmailConfig())
	authService := service.NewAuthService(authRepository, refreshTokenRepository, emailVerifyTokenRepository, smtpService)
	shortLinkService := service.NewShortLinkService(shortLinkRepository)

	log.Info("Registering Application Controller")
	authController := controller.NewAuthController(authService)
	shortLinkController := controller.NewShortLinkController(shortLinkService)
	redirectController := controller.NewRedirectController(shortLinkService)

	log.Info("Registered All API")
	return &APIContainer{
		AuthController:      authController,
		ShortLinkController: shortLinkController,
		RedirectController:  redirectController,
	}
}
