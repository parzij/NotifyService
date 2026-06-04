package models

import "time"

type Notification struct {
	ID           int64               `db:"id" json:"id"`
	UserID       int64               `db:"user_id" json:"user_id"`
	Recipient    string              `db:"recipient" json:"recipient"`
	Channel      NotificationChannel `db:"channel" json:"channel"`
	Message      string              `db:"message" json:"message"`
	Status       NotificationStatus  `db:"status" json:"status"`
	AttemptCount int16               `db:"attempt_count" json:"attempt_count"`
	ErrorMessage *string             `db:"error_message" json:"error_message,omitempty"`
	ExternalID   *string             `db:"external_id" json:"external_id,omitempty"`
	CreatedAt    time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time           `db:"updated_at" json:"updated_at"`
}
