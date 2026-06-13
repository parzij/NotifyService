// Package router provides HTTP router configuration
package router

import (
	"NotifyService/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler) {
	r.GET("/health", handlers.Health)

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	user := api.Group("/user")
	user.Use(authHandler.AuthMiddleware())
	user.PATCH("/phone", userHandler.UpdatePhone)
	user.PATCH("/telegram", userHandler.UpdateTelegram)
	user.GET("/profile", userHandler.GetProfile)
}
