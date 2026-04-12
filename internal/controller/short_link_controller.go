package controller

import (
	"oat431/fluffy-mouton/internal/payload/request"
	"oat431/fluffy-mouton/internal/payload/response"
	"oat431/fluffy-mouton/internal/service"
	"oat431/fluffy-mouton/pkg/common"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ShortLinkController struct {
	service service.ShortLinkService
}

func NewShortLinkController(service service.ShortLinkService) *ShortLinkController {
	return &ShortLinkController{service: service}
}

func (s *ShortLinkController) GetAllShortLinks(c fiber.Ctx) error {
	ownBy := c.Locals("auth_id").(uuid.UUID)
	shortLinkDTOs, err := s.service.GetAllLinks(c.Context(), ownBy)
	var res = common.ResponseDTO[[]response.ShortLinkDTO]{}
	if err != nil {
		res.Data = nil
		res.Status = common.ERROR
		res.Error = &common.ResponseDTOError{
			HttpCode:  fiber.StatusInternalServerError,
			ErrorCode: "INTERNAL_SERVER_ERROR",
			Message:   "Failed to retrieve short links",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}
	res.Data = &shortLinkDTOs
	res.Status = common.SUCCESS
	res.Error = nil

	return c.Status(fiber.StatusOK).JSON(res)
}

func (s *ShortLinkController) CreateRandomShortLink(c fiber.Ctx) error {
	req := c.Locals("payload").(*request.ShortLinkRequest)
	ownBy := c.Locals("auth_id").(uuid.UUID)

	shortLinkDTO, err := s.service.CreateRandomShortLink(c.Context(), req.Url, ownBy)
	var res = common.ResponseDTO[response.ShortLinkDTO]{}
	if err != nil {
		res.Data = nil
		res.Status = common.ERROR
		res.Error = &common.ResponseDTOError{
			HttpCode:  fiber.StatusInternalServerError,
			ErrorCode: "INTERNAL_SERVER_ERROR",
			Message:   "Failed to create short link",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}
	res.Data = shortLinkDTO
	res.Status = common.SUCCESS
	res.Error = nil

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (s *ShortLinkController) CreateCustomShortLink(c fiber.Ctx) error {
	req := c.Locals("payload").(*request.ShortLinkRequest)
	ownBy := c.Locals("auth_id").(uuid.UUID)

	shortLinkDTO, err := s.service.CreateCustomShortLink(c.Context(), req.Url, req.CustomName, ownBy)
	var res = common.ResponseDTO[response.ShortLinkDTO]{}
	if err != nil {
		res.Data = nil
		res.Status = common.ERROR
		res.Error = &common.ResponseDTOError{
			HttpCode:  fiber.StatusConflict,
			ErrorCode: "SHORT_LINK_ALREADY_EXISTS",
			Message:   "Custom short link already exists",
		}
		return c.Status(fiber.StatusConflict).JSON(res)
	}
	res.Data = shortLinkDTO
	res.Status = common.SUCCESS
	res.Error = nil

	return c.Status(fiber.StatusCreated).JSON(res)
}
