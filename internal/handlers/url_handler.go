package handlers

import (
	"ShortUrl/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShortenRequest struct {
	OriginalURL    string `json:"original_url" binding:"required"`
	CustomShortURL string `json:"custom_short_url"`
}

type URLHandler struct {
	URLService *services.URLService
}

func (h *URLHandler) CreateShortenUrl(c *gin.Context) {
	var req ShortenRequest

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userId, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	shortURl, err := h.URLService.SaveUrl(req.OriginalURL, req.CustomShortURL, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Create successful": shortURl})
}
