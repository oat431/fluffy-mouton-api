package controller_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"oat431/fluffy-mouton/internal/controller"
	"oat431/fluffy-mouton/internal/payload/response"
)

func TestRedirectController_ShortLinkRedirect(t *testing.T) {

	t.Run("redirects permanently on success", func(t *testing.T) {
		app := fiber.New()
		var gotCode, gotLinkType string
		svc := &mockShortLinkService{
			getLinkByCodeFunc: func(ctx context.Context, code string, linkType string) (*response.ShortLinkDTO, error) {
				gotCode = code
				gotLinkType = linkType
				return &response.ShortLinkDTO{
					OriginalLink: "https://example.com/long-url",
				}, nil
			},
		}

		ctrl, err := controller.NewRedirectController(svc)
		require.NoError(t, err)

		app.Get("/r/:code/:linkType", ctrl.ShortLinkRedirect)

		req := httptest.NewRequest("GET", "/r/abcde/r", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, "abcde", gotCode)
		assert.Equal(t, "RANDOM", gotLinkType)
		// Fiber 3's Redirect() currently defaults to 302/303 See Other if not chained correctly
		// but we check for any 3xx redirect
		assert.True(t, resp.StatusCode >= 300 && resp.StatusCode < 400)
		assert.Equal(t, "https://example.com/long-url", resp.Header.Get("Location"))
	})

	t.Run("returns 404 when link not found", func(t *testing.T) {
		app := fiber.New()
		svc := &mockShortLinkService{
			getLinkByCodeFunc: func(ctx context.Context, code string, linkType string) (*response.ShortLinkDTO, error) {
				return nil, errors.New("not found")
			},
		}

		ctrl, err := controller.NewRedirectController(svc)
		require.NoError(t, err)

		app.Get("/r/:code/:linkType", ctrl.ShortLinkRedirect)

		req := httptest.NewRequest("GET", "/r/notfound/c", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
}
