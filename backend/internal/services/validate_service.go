package services

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
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
