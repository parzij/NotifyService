// Package main for start point program
package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"NotifyService/internal/config"
	"NotifyService/internal/db"
	"NotifyService/internal/handlers"
	"NotifyService/internal/repository"
	"NotifyService/internal/router"
	"NotifyService/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	gin.SetMode(os.Getenv("GIN_MODE"))
	cfg := config.Load()

	dbPool, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	log.Println("✅ Database connection established successfully")

	port := cfg.HTTPport
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	userRepo := repository.NewUserRepo(dbPool)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, logger)
	authHandler := handlers.NewAuthHandler(authService)
	userService := services.NewUserService(userRepo, logger)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	router.Router(r, authHandler, userHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
