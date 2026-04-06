package middleware

import (
	"fmt"
	"oat431/fluffy-mouton/pkg/common"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

func Validate[T any](c fiber.Ctx) error {
	payload := new(T)

	if err := c.Bind().Body(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusBadRequest,
				ErrorCode: "BAD_REQUEST",
				Message:   "Invalid request body format",
			},
		})
	}

	if err := validate.Struct(payload); err != nil {
		var errorMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on tag '%s'", err.Field(), err.Tag())
			errorMessages = append(errorMessages, msg)
		}

		return c.Status(fiber.StatusBadRequest).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusBadRequest,
				ErrorCode: "VALIDATION_ERROR",
				Message:   strings.Join(errorMessages, ", "),
			},
		})
	}

	c.Locals("payload", payload)

	return c.Next()
}
