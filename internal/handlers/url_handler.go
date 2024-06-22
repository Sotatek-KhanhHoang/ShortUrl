package handlers

import (
	"ShortUrl/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShortenRequest struct {
	OriginalURL    string `json:"original_url" binding:"required"`
	CustomShortURL string `json:"custom_short_url" binding:"required"`
}

type URLHandler struct {
	URLService *services.URLService
}

func (h *URLHandler) CreateShortenUrl(c *gin.Context) {
	var req ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": userId})
		return
	}

	/*userId, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}*/

	shortURl, err := h.URLService.SaveUrl(req.OriginalURL, req.CustomShortURL, userId.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": userId})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Create successful": shortURl})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	shortenedURL := c.Param("shortenedURL")

	originalURL, err := h.URLService.GetOriginalURL(shortenedURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}
