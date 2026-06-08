package services

import (
	"NotifyService/internal/dto"
	"NotifyService/internal/repository"
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrContactPhoneAlreadyUsed    = errors.New("contact phone already used")
	ErrContactTelegramAlreadyUsed = errors.New("contact telegram already used")
)

func (s *AuthService) UpdatePhone(ctx context.Context, userID int64, req dto.PhoneRequest) (*dto.PhoneResponse, error) {
	phone := strings.TrimSpace(req.Phone)

	if phone == "" {
		return nil, ErrMissngFields
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

func (s *AuthService) UpdateTelegram(ctx context.Context, userID int64, req dto.TelegramRequest) (*dto.TelegramResponse, error) {
	telegram := strings.TrimSpace(req.Telegram)

	if telegram == "" {
		return nil, ErrMissngFields
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
