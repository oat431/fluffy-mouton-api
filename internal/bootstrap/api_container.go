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

func NewAPIContainer(db *sqlx.DB) (*APIContainer, error) {
	log.Info("Registering Application Repository")
	authRepository := repository.NewAuthRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)
	emailVerifyTokenRepository := repository.NewEmailVerifyTokenRepository(db)
	shortLinkRepository := repository.NewShortLinkRepository(db)

	log.Info("Registering Application Service")
	smtpService := service.NewSMTPService(config.GetEmailConfig())
	authService, err := service.NewAuthService(authRepository, refreshTokenRepository, emailVerifyTokenRepository, smtpService)
	if err != nil {
		return nil, err
	}
	shortLinkService, err := service.NewShortLinkService(shortLinkRepository)
	if err != nil {
		return nil, err
	}

	log.Info("Registering Application Controller")
	authController, err := controller.NewAuthController(authService)
	if err != nil {
		return nil, err
	}
	shortLinkController, err := controller.NewShortLinkController(shortLinkService)
	if err != nil {
		return nil, err
	}
	redirectController, err := controller.NewRedirectController(shortLinkService)
	if err != nil {
		return nil, err
	}

	log.Info("Registered All API")
	return &APIContainer{
		AuthController:      authController,
		ShortLinkController: shortLinkController,
		RedirectController:  redirectController,
	}, nil
}
