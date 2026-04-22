package controller

import (
	"errors"

	"oat431/fluffy-mouton/internal/payload/response"
	"oat431/fluffy-mouton/internal/service"
	"oat431/fluffy-mouton/pkg/common"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type RedirectController struct {
	service.ShortLinkService
}

func NewRedirectController(shortLinkService service.ShortLinkService) (*RedirectController, error) {
	if shortLinkService == nil {
		return nil, errors.New("redirect controller: nil service")
	}
	return &RedirectController{ShortLinkService: shortLinkService}, nil
}

func (s *RedirectController) ShortLinkRedirect(c fiber.Ctx) error {
	code := c.Params("code")
	linkType := linkTypeToEnum(c.Params("linkType"))

	shortLinkDTO, err := s.GetLinkByCode(c.Context(), code, linkType)
	if err != nil {
		res := common.ResponseDTO[response.ShortLinkDTO]{
			Data:   nil,
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusNotFound,
				ErrorCode: "SHORT_LINK_NOT_FOUND",
				Message:   "Short link not found",
			},
		}
		return c.Status(fiber.StatusNotFound).JSON(res)
	}
	log.Info("Redirecting to: ", shortLinkDTO.OriginalLink)
	return c.Status(fiber.StatusMovedPermanently).Redirect().To(shortLinkDTO.OriginalLink)
}

func linkTypeToEnum(linkType string) string {
	switch linkType {
	case "c":
		return "CUSTOM"
	case "r":
		return "RANDOM"
	default:
		return "r"
	}
}
