package repository

import (
	"context"
	"oat431/fluffy-mouton/internal/model"

	"github.com/jmoiron/sqlx"
)

type refreshTokenRepository struct {
	db *sqlx.DB
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, refreshToken model.RefreshToken) error
	Revoke(ctx context.Context, token string) error
	GetByToken(ctx context.Context, token string) (*model.RefreshToken, error)
}

func NewRefreshTokenRepository(db *sqlx.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Save(ctx context.Context, refreshToken model.RefreshToken) error {
	query := `INSERT INTO tb_refresh_tokens (
				id, created_at, updated_at, deleted_at, auth_id, token, expires_at, revoked
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query,
		refreshToken.ID,
		refreshToken.CreatedAt,
		refreshToken.UpdatedAt,
		refreshToken.DeletedAt,
		refreshToken.AuthID,
		refreshToken.Token,
		refreshToken.ExpiresAt,
		refreshToken.Revoked)
	return err
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, token string) error {
	query := `UPDATE tb_refresh_tokens SET revoked = true, updated_at = NOW() WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *refreshTokenRepository) GetByToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, auth_id, token, expires_at, revoked FROM tb_refresh_tokens WHERE token = $1 AND revoked = false AND deleted_at IS NULL`
	var refreshToken model.RefreshToken
	err := r.db.GetContext(ctx, &refreshToken, query, token)
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}
