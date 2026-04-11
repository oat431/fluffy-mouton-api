package router

import (
	"oat431/fluffy-mouton/internal/controller"

	"github.com/gofiber/fiber/v3"
)

func RegisterRedirectRoutes(router fiber.Router, controller *controller.RedirectController) {
	router.Get("/:linkType/:code", controller.ShortLinkRedirect)
}
