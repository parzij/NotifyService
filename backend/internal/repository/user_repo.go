package repository

import (
	"NotifyService/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdatePhone(ctx context.Context, userId int64, phone string) (*models.User, error)
	UpdateTelegramChatID(ctx context.Context, userId int64, telegramChatID string) (*models.User, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash, telegram_chat_id, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at;
	`

	return r.db.QueryRow(ctx, query,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.TelegramChatID,
		user.Phone,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, telegram_chat_id, phone, created_at, updated_at
		FROM users
		WHERE email = $1;
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.TelegramChatID,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepo) UpdatePhone(ctx context.Context, userId int64, phone string) (*models.User, error) {
	query := `
		UPDATE users 
		SET phone = $1,
		    updated_at = NOW()
		WHERE id = $2
		RETURNING id, full_name, email, password_hash, telegram_chat_id, phone, created_at, updated_at;    
	`

	user := &models.User{}

	err := r.db.QueryRow(ctx, query, phone, userId).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.TelegramChatID,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepo) UpdateTelegramChatID(ctx context.Context, userId int64, telegramChatID string) (*models.User, error) {
	query := `
          UPDATE users
          SET telegram_chat_id = $1,
              updated_at = NOW()
          WHERE id = $2
          RETURNING id, full_name, email, password_hash, telegram_chat_id, phone, created_at, updated_at;
	`

	user := &models.User{}

	err := r.db.QueryRow(ctx, query, telegramChatID, userId).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.TelegramChatID,
		&user.Phone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
