package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/mail"
	"strings"
	"time"

	"canhbuomxanh/backend/internal/config"
	"canhbuomxanh/backend/internal/platform/security"
)

type Service struct {
	repo *repository
	cfg  *config.Config
}

func NewService(cfg *config.Config, db *sql.DB) *Service {
	return &Service{
		repo: newRepository(db),
		cfg:  cfg,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) error {
	if req.Email == "" || req.Password == "" || req.FullName == "" || req.Role == "" {
		return errors.New("missing required fields")
	}
	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.New("invalid email format")
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	exists, err := s.repo.findUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists != nil {
		return errors.New("email already exists")
	}

	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return err
	}

	return s.repo.createUser(ctx, req, hash)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	user, err := s.repo.findUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !security.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	expiresAt := time.Now().Add(time.Duration(s.cfg.JWTExpiresHours) * time.Hour)
	token, err := security.GenerateJWT(user.ID, user.Role, s.cfg.JWTSecret, expiresAt)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(expiresAt).Seconds()),
	}, nil
}
