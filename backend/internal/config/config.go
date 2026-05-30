package config

import (
	"log"
	"os"
)

type Config struct {
	DBdsn     string
	HTTPport  string
	JWTSecret string
}

func Load() Config {
	httpPort := os.Getenv("APP_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSL := os.Getenv("DB_SSL")
	jwtSecret := os.Getenv("JWT_SECRET")

	if dbUser == "" || dbName == "" || dbHost == "" || dbPort == "" || dbSSL == "" {
		log.Fatal("DB_USER or DB_NAME or DB_HOST or DB_PORT or DB_SSL must be set")
	}

	if httpPort == "" {
		log.Fatal("APP port must be set")
	}

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	dsn := "postgres://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=" + dbSSL

	return Config{
		HTTPport:  httpPort,
		DBdsn:     dsn,
		JWTSecret: jwtSecret,
	}
}
