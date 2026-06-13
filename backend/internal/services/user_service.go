package services

import (
	"NotifyService/internal/dto"
	"NotifyService/internal/models"
	"NotifyService/internal/repository"
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserNotFound               = errors.New("user not found")
	ErrContactPhoneAlreadyUsed    = errors.New("contact phone already used")
	ErrContactTelegramAlreadyUsed = errors.New("contact telegram already used")
)

type UserService struct {
	repo   repository.UserRepository
	logger *slog.Logger
}

func NewUserService(repo repository.UserRepository, logger *slog.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		s.logger.Error(
			"failed to get user profile",
			slog.String("error", err.Error()),
			slog.Int64("user_id", userID),
		)

		return nil, ErrInternalServer
	}

	return user, nil
}

func (s *UserService) UpdatePhone(ctx context.Context, userID int64, req dto.PhoneRequest) (*dto.PhoneResponse, error) {
	phone := strings.TrimSpace(req.Phone)
	if phone == "" {
		return nil, ErrMissngFields
	}

	phone, err := ValidatePhone(phone)
	if err != nil {
		return nil, ErrInvalidPhone
	}

	user, err := s.repo.UpdatePhone(ctx, userID, phone)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrContactPhoneAlreadyUsed
		}

		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}

		s.logger.Error(
			"failed to update phone",
			slog.String("error", err.Error()),
			slog.Int64("user_id", userID),
		)
		return nil, ErrInternalServer
	}

	return &dto.PhoneResponse{User: *user}, nil
}

func (s *UserService) UpdateTelegram(ctx context.Context, userID int64, req dto.TelegramRequest) (*dto.TelegramResponse, error) {
	telegram := strings.TrimSpace(req.Telegram)
	if telegram == "" {
		return nil, ErrMissngFields
	}

	telegram, err := ValidateTelegram(telegram)
	if err != nil {
		return nil, ErrTelegramInvalidName
	}

	user, err := s.repo.UpdateTelegram(ctx, userID, telegram)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrContactTelegramAlreadyUsed
		}

		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}

		s.logger.Error(
			"failed to update telegram",
			slog.String("error", err.Error()),
			slog.Int64("user_id", userID),
		)

		return nil, ErrInternalServer
	}

	return &dto.TelegramResponse{User: *user}, nil
}
