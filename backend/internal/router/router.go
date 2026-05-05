package router

import (
    "github.com/gin-gonic/gin"
    "NotifyService/internal/handlers"
)

func Router(r *gin.Engine) {
	r.GET("/health", handlers.Health)

//	api := r.Group("/api/v1") {
		// notify
		// notification
//	}
}

