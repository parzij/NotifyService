-- +goose Up
-- +goose StatementBegin

ALTER TABLE users
RENAME COLUMN telegram_chat_id TO telegram;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE users
RENAME COLUMN telegram TO telegram_chat_id;

-- +goose StatementEnd