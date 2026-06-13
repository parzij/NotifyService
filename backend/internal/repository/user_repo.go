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
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	UpdatePhone(ctx context.Context, userID int64, phone string) (*models.User, error)
	UpdateTelegram(ctx context.Context, userID int64, telegram string) (*models.User, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash, telegram, phone)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at;
	`

	return r.db.QueryRow(ctx, query,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.Telegram,
		user.Phone,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, telegram, phone, created_at, updated_at
		FROM users
		WHERE email = $1;
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Telegram,
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

func (r *userRepo) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, telegram, phone, created_at, updated_at
		FROM users
		WHERE id = $1;
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Telegram,
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

func (r *userRepo) UpdatePhone(ctx context.Context, userID int64, phone string) (*models.User, error) {
	query := `
		UPDATE users 
		SET phone = $1,
		    updated_at = NOW()
		WHERE id = $2
		RETURNING id, full_name, email, password_hash, telegram, phone, created_at, updated_at;    
	`

	user := &models.User{}

	err := r.db.QueryRow(ctx, query, phone, userID).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Telegram,
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

func (r *userRepo) UpdateTelegram(ctx context.Context, userID int64, telegram string) (*models.User, error) {
	query := `
          UPDATE users
          SET telegram = $1,
              updated_at = NOW()
          WHERE id = $2
          RETURNING id, full_name, email, password_hash, telegram, phone, created_at, updated_at;
	`

	user := &models.User{}

	err := r.db.QueryRow(ctx, query, telegram, userID).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Telegram,
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
