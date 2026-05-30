// Package router provides HTTP router configuration
package router

import (
	"NotifyService/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine, authHandler *handlers.AuthHandler) {
	r.GET("/health", handlers.Health)

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
}
