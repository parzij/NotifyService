package handlers

import (
	"NotifyService/internal/dto"
	"NotifyService/internal/middleware"
	"NotifyService/internal/repository"
	"NotifyService/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) UpdatePhone(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not found in token"})
		return
	}
	var req dto.PhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.authService.UpdatePhone(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMissngFields):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrContactAlreadyUsed), errors.Is(err, services.ErrContactPhoneAlreadyUsed):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, repository.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": services.ErrInternalServer.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) UpdateTelegram(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not found in token"})
		return
	}

	var req dto.TelegramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.authService.UpdateTelegram(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMissngFields):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrContactAlreadyUsed), errors.Is(err, services.ErrContactTelegramAlreadyUsed):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, repository.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": services.ErrInternalServer.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, resp)
}
