-- +goose Up
-- +goose StatementBegin

CREATE TYPE notification_channel AS ENUM ('email', 'telegram', 'sms');
CREATE TYPE notification_status AS ENUM ('pending', 'sent', 'failed');

CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    telegram_chat_id VARCHAR(100) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    recipient VARCHAR(255) NOT NULL,
    channel notification_channel NOT NULL,
    message TEXT NOT NULL,
    status notification_status NOT NULL DEFAULT 'pending',
    attempt_count SMALLINT NOT NULL DEFAULT 0 CHECK (attempt_count >= 0),
    error_message TEXT,
    external_id VARCHAR(255) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_channel_status ON notifications(channel, status);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_notifications_channel_status;
DROP INDEX IF EXISTS idx_notifications_created_at;
DROP INDEX IF EXISTS idx_notifications_status;
DROP INDEX IF EXISTS idx_notifications_user_id;

DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS notification_status;
DROP TYPE IF EXISTS notification_channel;

-- +goose StatementEnd