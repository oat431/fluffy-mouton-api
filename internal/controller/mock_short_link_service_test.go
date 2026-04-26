package controller_test

import (
	"context"
	"oat431/fluffy-mouton/internal/payload/response"
	"oat431/fluffy-mouton/internal/service"

	"github.com/google/uuid"
)

var _ service.ShortLinkService = (*mockShortLinkService)(nil)

type mockShortLinkService struct {
	getAllLinksFunc           func(ctx context.Context, ownBy uuid.UUID) ([]response.ShortLinkDTO, error)
	getLinkByCodeFunc         func(ctx context.Context, code string, linkType string) (*response.ShortLinkDTO, error)
	createRandomShortLinkFunc func(ctx context.Context, originalURL string, ownBy uuid.UUID) (*response.ShortLinkDTO, error)
	createCustomShortLinkFunc func(ctx context.Context, originalURL string, customCode string, ownBy uuid.UUID) (*response.ShortLinkDTO, error)
	updateShortLinkFunc       func(ctx context.Context, id string, url string, ownBy uuid.UUID) (*response.ShortLinkDTO, error)
	deleteShortLinkFunc       func(ctx context.Context, id string, ownBy uuid.UUID) error
}

func (m *mockShortLinkService) GetAllLinks(ctx context.Context, ownBy uuid.UUID) ([]response.ShortLinkDTO, error) {
	if m.getAllLinksFunc != nil {
		return m.getAllLinksFunc(ctx, ownBy)
	}
	return nil, nil
}

func (m *mockShortLinkService) GetLinkByCode(ctx context.Context, code string, linkType string) (*response.ShortLinkDTO, error) {
	if m.getLinkByCodeFunc != nil {
		return m.getLinkByCodeFunc(ctx, code, linkType)
	}
	return nil, nil
}

func (m *mockShortLinkService) CreateRandomShortLink(ctx context.Context, originalURL string, ownBy uuid.UUID) (*response.ShortLinkDTO, error) {
	if m.createRandomShortLinkFunc != nil {
		return m.createRandomShortLinkFunc(ctx, originalURL, ownBy)
	}
	return nil, nil
}

func (m *mockShortLinkService) CreateCustomShortLink(ctx context.Context, originalURL string, customCode string, ownBy uuid.UUID) (*response.ShortLinkDTO, error) {
	if m.createCustomShortLinkFunc != nil {
		return m.createCustomShortLinkFunc(ctx, originalURL, customCode, ownBy)
	}
	return nil, nil
}

func (m *mockShortLinkService) UpdateShortLink(ctx context.Context, id string, url string, ownBy uuid.UUID) (*response.ShortLinkDTO, error) {
	if m.updateShortLinkFunc != nil {
		return m.updateShortLinkFunc(ctx, id, url, ownBy)
	}
	return nil, nil
}

func (m *mockShortLinkService) DeleteShortLink(ctx context.Context, id string, ownBy uuid.UUID) error {
	if m.deleteShortLinkFunc != nil {
		return m.deleteShortLinkFunc(ctx, id, ownBy)
	}
	return nil
}
