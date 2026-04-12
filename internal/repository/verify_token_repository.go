package repository

import (
	"context"
	"oat431/fluffy-mouton/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type emailVerifyTokenRepository struct {
	db *sqlx.DB
}

type EmailVerifyTokenRepository interface {
	Save(ctx context.Context, token model.EmailVerifyToken) error
	FindByToken(ctx context.Context, token string) (*model.EmailVerifyToken, error)
	DeleteByAuthID(ctx context.Context, authID uuid.UUID) error
}

func NewEmailVerifyTokenRepository(db *sqlx.DB) EmailVerifyTokenRepository {
	return &emailVerifyTokenRepository{db: db}
}

func (r *emailVerifyTokenRepository) Save(ctx context.Context, token model.EmailVerifyToken) error {
	query := `INSERT INTO tb_verify_tokens (
				id, created_at, updated_at, deleted_at, auth_id, token, expires_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.CreatedAt,
		token.UpdatedAt,
		token.DeletedAt,
		token.AuthID,
		token.Token,
		token.ExpiresAt,
	)
	return err
}

func (r *emailVerifyTokenRepository) FindByToken(ctx context.Context, token string) (*model.EmailVerifyToken, error) {
	query := `SELECT id, created_at, updated_at, deleted_at, auth_id, token, expires_at
			  FROM tb_verify_tokens
			  WHERE token = $1 AND deleted_at IS NULL`
	var result model.EmailVerifyToken
	err := r.db.GetContext(ctx, &result, query, token)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *emailVerifyTokenRepository) DeleteByAuthID(ctx context.Context, authID uuid.UUID) error {
	query := `DELETE FROM tb_verify_tokens WHERE auth_id = $1`
	_, err := r.db.ExecContext(ctx, query, authID)
	return err
}
