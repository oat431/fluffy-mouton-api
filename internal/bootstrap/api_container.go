package bootstrap

import (
	"oat431/fluffy-mouton/internal/controller"
	"oat431/fluffy-mouton/internal/repository"
	"oat431/fluffy-mouton/internal/service"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jmoiron/sqlx"
)

type APIContainer struct {
	ShortLinkController *controller.ShortLinkController
	RedirectController  *controller.RedirectController
}

func NewAPIContainer(db *sqlx.DB) (*APIContainer, error) {
	log.Info("Registering Application Repository")
	shortLinkRepository := repository.NewShortLinkRepository(db)

	log.Info("Registering Application Service")
	shortLinkService, err := service.NewShortLinkService(shortLinkRepository)
	if err != nil {
		return nil, err
	}

	log.Info("Registering Application Controller")
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
		ShortLinkController: shortLinkController,
		RedirectController:  redirectController,
	}, nil
}
