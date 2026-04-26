package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"oat431/fluffy-mouton/internal/model"
	"oat431/fluffy-mouton/internal/service"
	"oat431/fluffy-mouton/pkg/common"
)

func TestShortLinkService_CreateCustomShortLink(t *testing.T) {
	ctx := context.Background()
	ownBy := uuid.New()
	originalURL := "https://example.com"
	customName := "mycustom"

	t.Run("succeeds when custom name is unique", func(t *testing.T) {
		repo := &mockShortLinkRepo{
			getLinkByShortCodeFunc: func(ctx context.Context, code string, linkType string) (*model.ShortLink, error) {
				return nil, nil // Not found, unique name
			},
			createShortLinkFunc: func(ctx context.Context, url string, shortUrl string, linkType string, ownBy uuid.UUID) (*model.ShortLink, error) {
				return &model.ShortLink{
					TargetURL: url,
					ShortURL:  shortUrl,
					Type:      model.LinkType(linkType),
				}, nil
			},
		}

		svc, err := service.NewShortLinkService(repo)
		require.NoError(t, err)

		got, err := svc.CreateCustomShortLink(ctx, originalURL, customName, ownBy)
		require.NoError(t, err)
		assert.Equal(t, customName, got.ShortLink)
		assert.Equal(t, originalURL, got.OriginalLink)
		assert.Equal(t, "CUSTOM", got.LinkType)
	})

	t.Run("returns error when custom name already exists", func(t *testing.T) {
		repo := &mockShortLinkRepo{
			getLinkByShortCodeFunc: func(ctx context.Context, code string, linkType string) (*model.ShortLink, error) {
				return &model.ShortLink{ShortURL: code}, nil // Found existing
			},
		}

		svc, err := service.NewShortLinkService(repo)
		require.NoError(t, err)

		_, err = svc.CreateCustomShortLink(ctx, originalURL, customName, ownBy)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestShortLinkService_GetLinkByCode(t *testing.T) {
	ctx := context.Background()
	code := "abcde"
	linkType := "RANDOM"
	id := uuid.New()

	t.Run("succeeds and updates view count", func(t *testing.T) {
		updateCalled := false
		repo := &mockShortLinkRepo{
			getLinkByShortCodeFunc: func(ctx context.Context, c string, lt string) (*model.ShortLink, error) {
				return &model.ShortLink{
					BaseEntity: common.BaseEntity{ID: id},
					TargetURL:  "https://example.com",
					ShortURL:   c,
					Type:       model.LinkType(lt),
					View:       5,
				}, nil
			},
			updateViewCountFunc: func(ctx context.Context, linkID string, view int) error {
				assert.Equal(t, id.String(), linkID)
				assert.Equal(t, 6, view)
				updateCalled = true
				return nil
			},
		}

		svc, err := service.NewShortLinkService(repo)
		require.NoError(t, err)

		got, err := svc.GetLinkByCode(ctx, code, linkType)
		require.NoError(t, err)
		assert.Equal(t, code, got.ShortLink)
		assert.Equal(t, "https://example.com", got.OriginalLink)
		assert.True(t, updateCalled, "UpdateViewCount should have been called")
	})

	t.Run("returns error when link not found", func(t *testing.T) {
		repo := &mockShortLinkRepo{
			getLinkByShortCodeFunc: func(ctx context.Context, c string, lt string) (*model.ShortLink, error) {
				return nil, errors.New("not found")
			},
		}

		svc, err := service.NewShortLinkService(repo)
		require.NoError(t, err)

		_, err = svc.GetLinkByCode(ctx, code, linkType)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
