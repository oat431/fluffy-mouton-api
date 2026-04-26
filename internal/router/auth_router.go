package router

import (
	"oat431/fluffy-mouton/internal/controller"
	"oat431/fluffy-mouton/internal/middleware"
	"oat431/fluffy-mouton/internal/payload/request"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(router fiber.Router, controller *controller.AuthController) {
	route := router.Group("/auth")

	route.Post("/register",
		middleware.Validate[request.RegisterRequest],
		controller.RegisterNewUser,
	)

	route.Post("/login",
		middleware.Validate[request.LoginRequest],
		controller.Login,
	)

	route.Post("/revoke",
		controller.RevokeAccess,
	)

	route.Get("/detail",
		middleware.JWTMiddleware,
		controller.GetUserDetails,
	)

	route.Get("/verify-email",
		controller.VerifyEmail,
	)
}
