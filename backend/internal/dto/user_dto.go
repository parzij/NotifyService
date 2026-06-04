package dto

import "NotifyService/internal/models"

type RegisterRequest struct {
	FullName       string  `json:"full_name"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	TelegramChatID *string `json:"telegram_chat_id"`
	Phone          *string `json:"phone"`
}

type RegisterResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

type PhoneRequest struct {
	Phone string `json:"phone"`
}

type PhoneResponse struct {
	User models.User `json:"user"`
}

type TelegramRequest struct {
	TelegramChatID string `json:"telegram_chat_id"`
}

type TelegramResponse struct {
	User models.User `json:"user"`
}
