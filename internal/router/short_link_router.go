package router

import (
	"oat431/fluffy-mouton/internal/controller"
	"oat431/fluffy-mouton/internal/middleware"
	"oat431/fluffy-mouton/internal/payload/request"

	"github.com/gofiber/fiber/v3"
)

func RegisterShortLinkRoutes(router fiber.Router, controller *controller.ShortLinkController) {
	route := router.Group("/short")

	route.Get("/",
		middleware.OAuthMiddleware,
		controller.GetAllShortLinks,
	)
	route.Post("/random",
		middleware.Validate[request.ShortLinkRequest],
		middleware.OAuthMiddleware,
		controller.CreateRandomShortLink,
	)
	route.Post("/custom",
		middleware.Validate[request.ShortLinkRequest],
		middleware.OAuthMiddleware,
		controller.CreateCustomShortLink,
	)
	route.Put("/:id",
		middleware.Validate[request.UpdateShortLinkRequest],
		middleware.OAuthMiddleware,
		controller.UpdateShortLink,
	)
	route.Delete("/:id",
		middleware.OAuthMiddleware,
		controller.DeleteShortLink,
	)
}
