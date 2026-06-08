package services

import (
	"NotifyService/internal/dto"
	"NotifyService/internal/models"
	"NotifyService/internal/repository"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrMissngFields       = errors.New("error missing fields")
	ErrWeakPassword       = errors.New("password must be at least 8 characters long")
	ErrInternalServer     = errors.New("internal server error")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrContactAlreadyUsed = errors.New("phone or telegram chat id already used")
)

type AuthService struct {
	repo      repository.UserRepository
	jwtSecret []byte
	logger    *slog.Logger
}

func NewAuthService(repo repository.UserRepository, jwtSecret string, logger *slog.Logger) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		logger:    logger,
	}
}

func (s *AuthService) generateToken(userID int64, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *AuthService) VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return 0, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id not found in token")
	}

	return int64(userID), nil
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	req.FullName = strings.TrimSpace(req.FullName)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.FullName == "" || req.Email == "" || req.Password == "" {
		return nil, ErrMissngFields
	}

	if len(req.Password) < 8 {
		return nil, ErrWeakPassword
	}

	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", slog.String("error", err.Error()))
		return nil, ErrInternalServer
	}

	user := &models.User{
		FullName:     req.FullName,
		Email:        strings.ToLower(req.Email),
		PasswordHash: passwordHash,
		Telegram:     req.Telegram,
		Phone:        req.Phone,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrUserAlreadyExists
		}

		s.logger.Error("failed to create user", slog.String("error", err.Error()), slog.String("email", user.Email))
		return nil, ErrInternalServer
	}

	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		s.logger.Error("failed to generate token", slog.String("error", err.Error()), slog.String("email", user.Email))
		return nil, ErrInternalServer
	}

	s.logger.Info("user successfully registered", slog.Int64("user_id", user.ID), slog.String("email", user.Email))
	return &dto.RegisterResponse{User: *user, Token: token}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		return nil, ErrMissngFields
	}

	user, err := s.repo.GetUserByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		s.logger.Error("failed to find user", slog.String("error", err.Error()), slog.String("email", req.Email))
		return nil, ErrInternalServer
	}

	if err := CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		s.logger.Error("failed to generate token", slog.String("error", err.Error()), slog.String("email", user.Email))
		return nil, ErrInternalServer
	}

	s.logger.Info("user successfully logged in", slog.Int64("user_id", user.ID), slog.String("email", user.Email))
	return &dto.LoginResponse{User: *user, Token: token}, nil
}
