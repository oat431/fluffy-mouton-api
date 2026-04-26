package service_test

import (
	"context"
	"oat431/fluffy-mouton/internal/model"
	"oat431/fluffy-mouton/internal/repository"

	"github.com/google/uuid"
)

var _ repository.ShortLinkRepository = (*mockShortLinkRepo)(nil)

type mockShortLinkRepo struct {
	getAllShortLinkFunc    func(ctx context.Context, ownBy uuid.UUID) ([]model.ShortLink, error)
	getLinkByShortCodeFunc func(ctx context.Context, code string, linkType string) (*model.ShortLink, error)
	createShortLinkFunc    func(ctx context.Context, url string, shortUrl string, linkType string, ownBy uuid.UUID) (*model.ShortLink, error)
	updateViewCountFunc    func(ctx context.Context, id string, view int) error
	updateShortLinkURLFunc func(ctx context.Context, id string, url string, ownBy uuid.UUID) (*model.ShortLink, error)
	deleteShortLinkFunc    func(ctx context.Context, id string, ownBy uuid.UUID) error
}

func (m *mockShortLinkRepo) GetAllShortLink(ctx context.Context, ownBy uuid.UUID) ([]model.ShortLink, error) {
	if m.getAllShortLinkFunc != nil {
		return m.getAllShortLinkFunc(ctx, ownBy)
	}
	return nil, nil
}

func (m *mockShortLinkRepo) GetLinkByShortCode(ctx context.Context, code string, linkType string) (*model.ShortLink, error) {
	if m.getLinkByShortCodeFunc != nil {
		return m.getLinkByShortCodeFunc(ctx, code, linkType)
	}
	return nil, nil
}

func (m *mockShortLinkRepo) CreateShortLink(ctx context.Context, url string, shortUrl string, linkType string, ownBy uuid.UUID) (*model.ShortLink, error) {
	if m.createShortLinkFunc != nil {
		return m.createShortLinkFunc(ctx, url, shortUrl, linkType, ownBy)
	}
	return nil, nil
}

func (m *mockShortLinkRepo) UpdateViewCount(ctx context.Context, id string, view int) error {
	if m.updateViewCountFunc != nil {
		return m.updateViewCountFunc(ctx, id, view)
	}
	return nil
}

func (m *mockShortLinkRepo) UpdateShortLinkURL(ctx context.Context, id string, url string, ownBy uuid.UUID) (*model.ShortLink, error) {
	if m.updateShortLinkURLFunc != nil {
		return m.updateShortLinkURLFunc(ctx, id, url, ownBy)
	}
	return nil, nil
}

func (m *mockShortLinkRepo) DeleteShortLink(ctx context.Context, id string, ownBy uuid.UUID) error {
	if m.deleteShortLinkFunc != nil {
		return m.deleteShortLinkFunc(ctx, id, ownBy)
	}
	return nil
}
