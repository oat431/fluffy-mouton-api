package repository

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"oat431/fluffy-mouton/internal/model"
)

type ShortLinkRepository interface {
	GetAllShortLink(ctx context.Context, ownBy uuid.UUID) ([]model.ShortLink, error)
	GetLinkByShortCode(ctx context.Context, code string, linkType string) (*model.ShortLink, error)
	CreateShortLink(ctx context.Context, url string, shortUrl string, linkType string, ownBy uuid.UUID) (*model.ShortLink, error)
	UpdateViewCount(ctx context.Context, id string, view int) error
}

type shortLinkRepository struct {
	db *sqlx.DB
}

func NewShortLinkRepository(db *sqlx.DB) ShortLinkRepository {
	return &shortLinkRepository{db: db}
}

func (s shortLinkRepository) GetAllShortLink(ctx context.Context, ownBy uuid.UUID) ([]model.ShortLink, error) {
	query := "SELECT id, target_url, short_url, type, created_at FROM tb_short_links WHERE own_by = $1"
	rows, err := s.db.QueryContext(ctx, query, ownBy)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shortLinks []model.ShortLink
	for rows.Next() {
		var sl model.ShortLink
		if err := rows.Scan(&sl.ID, &sl.TargetURL, &sl.ShortURL, &sl.Type, &sl.CreatedAt); err != nil {
			return nil, err
		}
		shortLinks = append(shortLinks, sl)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return shortLinks, nil
}

func (s shortLinkRepository) GetLinkByShortCode(ctx context.Context, code string, linkType string) (*model.ShortLink, error) {
	query := "SELECT id, target_url, short_url, type, created_at FROM tb_short_links WHERE short_url = $1 AND type = $2 "
	row := s.db.QueryRowContext(ctx, query, code, linkType)

	var sl model.ShortLink
	err := row.Scan(&sl.ID, &sl.TargetURL, &sl.ShortURL, &sl.Type, &sl.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &sl, nil
}

func (s shortLinkRepository) CreateShortLink(ctx context.Context, url string, shortUrl string, linkType string, ownBy uuid.UUID) (*model.ShortLink, error) {
	query := "INSERT INTO tb_short_links (id,target_url, short_url, type, created_at, own_by) VALUES ($1, $2, $3,$4,$5,$6) RETURNING id, target_url, short_url, type, created_at, own_by"
	id := uuid.New()
	var sl model.ShortLink
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
		url,
		shortUrl,
		linkType,
		time.Now(),
		ownBy,
	).Scan(&sl.ID, &sl.TargetURL, &sl.ShortURL, &sl.Type, &sl.CreatedAt, &sl.OwnBy)
	if err != nil {
		log.Error("Error inserting short link: ", err)
		return nil, err
	}
	return &sl, nil
}

func (s shortLinkRepository) UpdateViewCount(ctx context.Context, id string, view int) error {
	query := "UPDATE tb_short_links SET view = $1 WHERE id = $2"
	_, err := s.db.ExecContext(ctx, query, view, id)
	return err
}
