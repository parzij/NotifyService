package services

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPhone        = errors.New("invalid phone number")
	ErrTelegramInvalidName = errors.New("telegram invalid name")
)

func HashPassword(plain string) (string, error) {
	if plain == "" {
		return "", errors.New("password is empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}
	return string(hash), nil
}

func CheckPasswordHash(plain, hash string) error {
	if plain == "" || hash == "" {
		return errors.New("password or hash is empty")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)); err != nil {
		return fmt.Errorf("bcrypt.CompareHashAndPassword: %w", err)
	}

	return nil
}

func ValidatePhone(phone string) (string, error) {
	if len(phone) != 12 || phone[0] != '+' || phone[1] != '7' {
		return "", ErrInvalidPhone
	}

	for i := 2; i < len(phone); i++ {
		if phone[i] < '0' || phone[i] > '9' {
			return "", ErrInvalidPhone
		}
	}
	return phone, nil
}

func ValidateTelegram(telegram string) (string, error) {
	if len(telegram) < 6 || len(telegram) > 33 {
		return "", ErrTelegramInvalidName
	}

	if telegram[0] != '@' {
		return "", ErrTelegramInvalidName
	}

	first := telegram[1]
	if !((first >= 'a' && first <= 'z') ||
		(first >= 'A' && first <= 'Z')) {
		return "", ErrTelegramInvalidName
	}

	for i := 1; i < len(telegram); i++ {
		ch := telegram[i]

		isDigit := ch >= '0' && ch <= '9'
		isLowerLetter := ch >= 'a' && ch <= 'z'
		isUpperLetter := ch >= 'A' && ch <= 'Z'
		isUnderscore := ch == '_'

		if !isDigit && !isLowerLetter && !isUpperLetter && !isUnderscore {
			return "", ErrTelegramInvalidName
		}
	}

	return telegram, nil
}
