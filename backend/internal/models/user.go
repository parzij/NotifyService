package models

import "time"

type User struct {
	ID           int64     `db:"id" json:"id"`
	FullName     string    `db:"full_name" json:"full_name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Telegram     *string   `db:"telegram" json:"telegram,omitempty"`
	Phone        *string   `db:"phone" json:"phone,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
