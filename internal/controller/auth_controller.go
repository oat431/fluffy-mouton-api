package controller

import (
	"errors"

	"oat431/fluffy-mouton/internal/payload/request"
	"oat431/fluffy-mouton/internal/payload/response"
	"oat431/fluffy-mouton/internal/service"
	"oat431/fluffy-mouton/pkg/common"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) (*AuthController, error) {
	if service == nil {
		return nil, errors.New("auth controller: nil service")
	}
	return &AuthController{service: service}, nil
}

func (auth *AuthController) RegisterNewUser(c fiber.Ctx) error {
	req := c.Locals("payload").(*request.RegisterRequest)
	authDto, err := auth.service.Register(c.Context(), *req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.ErrBadRequest.Code,
				ErrorCode: "REGISTER-01",
				Message:   err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(common.ResponseDTO[response.AuthResponse]{
		Status: common.SUCCESS,
		Data:   authDto,
	})
}

func (auth *AuthController) Login(c fiber.Ctx) error {
	req := c.Locals("payload").(*request.LoginRequest)
	tokenDto, err := auth.service.Login(c.Context(), *req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.ErrBadRequest.Code,
				ErrorCode: "LOGIN-01",
				Message:   "Invalid username or password",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.ResponseDTO[response.JWTResponse]{
		Status: common.SUCCESS,
		Data:   tokenDto,
	})
}

func (auth *AuthController) RevokeAccess(c fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.Bind().JSON(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.ErrBadRequest.Code,
				ErrorCode: "REVOKE-01",
				Message:   "Invalid request body",
			},
		})
	}

	err := auth.service.RevokeAccess(c.Context(), body.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.ErrInternalServerError.Code,
				ErrorCode: "REVOKE-02",
				Message:   err.Error(),
			},
		})
	}

	message := "Access revoked"
	return c.Status(fiber.StatusOK).JSON(common.ResponseDTO[string]{
		Status: common.SUCCESS,
		Data:   &message,
	})
}

func (auth *AuthController) GetUserDetails(c fiber.Ctx) error {
	authID := c.Locals("auth_id").(uuid.UUID)
	authDto, err := auth.service.GetUserDetails(c.Context(), authID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.ErrInternalServerError.Code,
				ErrorCode: "DETAIL-01",
				Message:   err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.ResponseDTO[response.AuthResponse]{
		Status: common.SUCCESS,
		Data:   authDto,
	})
}

func (auth *AuthController) VerifyEmail(c fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.ErrBadRequest.Code,
				ErrorCode: "VERIFY-01",
				Message:   "Missing verification token",
			},
		})
	}

	if err := auth.service.VerifyEmail(c.Context(), token); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.ErrBadRequest.Code,
				ErrorCode: "VERIFY-02",
				Message:   err.Error(),
			},
		})
	}

	message := "Email verified successfully"
	return c.Status(fiber.StatusOK).JSON(common.ResponseDTO[string]{
		Status: common.SUCCESS,
		Data:   &message,
	})
}
