package handlers

import (
	"NotifyService/internal/dto"
	"NotifyService/internal/middleware"
	"NotifyService/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return middleware.AuthRequired(h.authService)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMissngFields), errors.Is(err, services.ErrWeakPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": services.ErrInternalServer.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMissngFields):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": services.ErrInternalServer.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
