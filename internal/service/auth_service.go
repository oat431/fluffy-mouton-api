package service

import (
	"context"
	"errors"
	"oat431/fluffy-mouton/internal/model"
	"oat431/fluffy-mouton/internal/payload/request"
	"oat431/fluffy-mouton/internal/payload/response"
	"oat431/fluffy-mouton/internal/repository"
	"oat431/fluffy-mouton/pkg/common"
	"oat431/fluffy-mouton/pkg/utils"
	"time"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v3/log"
)

type authService struct {
	repo                 repository.AuthRepository
	refreshTokenRepo     repository.RefreshTokenRepository
	emailVerifyTokenRepo repository.EmailVerifyTokenRepository
	emailService         *SMTPService
}

type AuthService interface {
	Register(ctx context.Context, request request.RegisterRequest) (*response.AuthResponse, error)
	Login(ctx context.Context, request request.LoginRequest) (*response.JWTResponse, error)
	RevokeAccess(ctx context.Context, refreshToken string) error
	GetUserDetails(ctx context.Context, authID uuid.UUID) (*response.AuthResponse, error)
	VerifyEmail(ctx context.Context, token string) error
}

func NewAuthService(
	repo repository.AuthRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	emailVerifyTokenRepo repository.EmailVerifyTokenRepository,
	emailService *SMTPService,
) (AuthService, error) {
	if repo == nil {
		return nil, errors.New("auth service: nil auth repository")
	}
	if refreshTokenRepo == nil {
		return nil, errors.New("auth service: nil refresh token repository")
	}
	if emailVerifyTokenRepo == nil {
		return nil, errors.New("auth service: nil email verify token repository")
	}
	if emailService == nil {
		return nil, errors.New("auth service: nil email service")
	}

	return &authService{
		repo:                 repo,
		refreshTokenRepo:     refreshTokenRepo,
		emailVerifyTokenRepo: emailVerifyTokenRepo,
		emailService:         emailService,
	}, nil
}

func (s *authService) Register(ctx context.Context, request request.RegisterRequest) (*response.AuthResponse, error) {
	encryptedPassword, err := utils.EncryptPassword(request.Password)
	if err != nil {
		return nil, err
	}
	request.Password = encryptedPassword
	auth, err := s.repo.Register(ctx, request)
	if err != nil {
		return nil, err
	}

	// Generate email verification token
	tokenStr := utils.GenerateVerifyToken()
	if tokenStr == "" {
		return nil, errors.New("failed to generate verification token")
	}

	emailVerifyToken := model.EmailVerifyToken{
		BaseEntity: common.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		AuthID:    auth.ID,
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	if err := s.emailVerifyTokenRepo.Save(ctx, emailVerifyToken); err != nil {
		log.Error("failed to save email verify token: ", err.Error())
		return nil, err
	}

	// Send verification email (non-blocking on failure to keep registration flow intact)
	if err := s.emailService.SendVerificationEmail(auth.Email, tokenStr); err != nil {
		log.Error("failed to send verification email: ", err.Error())
	}

	return &response.AuthResponse{
		ID:         utils.UUIDToString(auth.ID),
		Username:   auth.Username,
		Email:      auth.Email,
		IsVerified: auth.IsVerified,
	}, nil
}

func (s *authService) Login(ctx context.Context, request request.LoginRequest) (*response.JWTResponse, error) {
	auth, err := s.repo.GetAuthByUsername(ctx, request.Username)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = utils.ComparePassword(auth.Password, request.Password)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	accessToken, err := utils.GenerateAccessToken(auth.ID)
	if err != nil {
		return nil, err
	}

	refreshTokenStr, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshToken := model.RefreshToken{
		BaseEntity: common.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		AuthID:    auth.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // 7 days
		Revoked:   false,
	}

	err = s.refreshTokenRepo.Save(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return &response.JWTResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (s *authService) RevokeAccess(ctx context.Context, refreshToken string) error {
	return s.refreshTokenRepo.Revoke(ctx, refreshToken)
}

func (s *authService) GetUserDetails(ctx context.Context, authID uuid.UUID) (*response.AuthResponse, error) {
	auth, err := s.repo.GetAuthByID(ctx, authID)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		ID:         utils.UUIDToString(auth.ID),
		Username:   auth.Username,
		Email:      auth.Email,
		IsVerified: auth.IsVerified,
	}, nil
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	emailVerifyToken, err := s.emailVerifyTokenRepo.FindByToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	if time.Now().After(emailVerifyToken.ExpiresAt) {
		return errors.New("verification token has expired")
	}

	if err := s.repo.MarkAsVerified(ctx, emailVerifyToken.AuthID); err != nil {
		return err
	}

	// Consume (delete) the token after successful verification
	if err := s.emailVerifyTokenRepo.DeleteByAuthID(ctx, emailVerifyToken.AuthID); err != nil {
		log.Error("failed to delete email verify token: ", err.Error())
	}

	return nil
}
