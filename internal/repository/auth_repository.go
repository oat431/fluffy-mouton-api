package repository

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"oat431/fluffy-mouton/internal/model"
	"oat431/fluffy-mouton/internal/payload/request"
	"oat431/fluffy-mouton/pkg/common"
)

type authRepository struct {
	db *sqlx.DB
}

type AuthRepository interface {
	Register(ctx context.Context, request request.RegisterRequest) (*model.Auth, error)
	GetAuthByUsername(ctx context.Context, username string) (*model.Auth, error)
	GetAuthByID(ctx context.Context, id uuid.UUID) (*model.Auth, error)
	GetAuthByEmail(ctx context.Context, email string) (*model.Auth, error)
	MarkAsVerified(ctx context.Context, authID uuid.UUID) error
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Register(ctx context.Context, request request.RegisterRequest) (*model.Auth, error) {
	query := `INSERT INTO tb_auth (
				id,
				created_at,
				updated_at,
				deleted_at,
				username,
				email,
				"password",
				is_verified
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	id := uuid.New()
	currentTime := time.Now()
	_, err := r.db.ExecContext(ctx, query, id, currentTime, currentTime, nil, request.Username, request.Email, request.Password, false)
	if err != nil {
		return nil, err
	}
	return &model.Auth{
		BaseEntity: common.BaseEntity{
			ID:        id,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: nil,
		},
		Username:   request.Username,
		Email:      request.Email,
		Password:   request.Password,
		IsVerified: false,
	}, nil
}

func (r *authRepository) GetAuthByUsername(ctx context.Context, username string) (*model.Auth, error) {
	query := `
		SELECT 
			id,
			created_at,
			updated_at,
			deleted_at,
			username,
			email,
			"password",
			is_verified
		FROM
			tb_auth
		WHERE
			username = $1`
	var auth model.Auth
	err := r.db.GetContext(ctx, &auth, query, username)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &auth, nil
}

func (r *authRepository) GetAuthByID(ctx context.Context, id uuid.UUID) (*model.Auth, error) {
	query := `
		SELECT 
			id,
			created_at,
			updated_at,
			deleted_at,
			username,
			email,
			"password",
			is_verified
		FROM
			tb_auth
		WHERE
			id = $1`
	var auth model.Auth
	err := r.db.GetContext(ctx, &auth, query, id)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &auth, nil
}

func (r *authRepository) GetAuthByEmail(ctx context.Context, email string) (*model.Auth, error) {
	return nil, nil
}

func (r *authRepository) MarkAsVerified(ctx context.Context, authID uuid.UUID) error {
	query := `UPDATE tb_auth SET is_verified = true, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, authID)
	return err
}
