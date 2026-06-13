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

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	value, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user is not authorized",
		})
		return
	}

	userID, ok := value.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user identifier",
		})
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *UserHandler) UpdatePhone(c *gin.Context) {
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

	resp, err := h.userService.UpdatePhone(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPhone):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid phone number",
			})
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

func (h *UserHandler) UpdateTelegram(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req dto.TelegramRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	response, err := h.userService.UpdateTelegram(
		c.Request.Context(),
		userID,
		req,
	)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMissngFields):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "telegram is required",
			})

		case errors.Is(err, services.ErrTelegramInvalidName):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid telegram username",
			})

		case errors.Is(err, services.ErrContactTelegramAlreadyUsed):
			c.JSON(http.StatusConflict, gin.H{
				"error": "telegram already used",
			})

		case errors.Is(err, services.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, response)
}
