package handlers

import (
	"ShortUrl/internal/config"
	"ShortUrl/internal/services"

	"net/http"
	"strings"

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

func (h *URLHandler) CreateRandomUrl(c *gin.Context) {
	var req ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.HasPrefix(req.CustomShortURL, "http://") || strings.HasPrefix(req.CustomShortURL, "https://") {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Invalid Url"})
		return
	}

	randomURL, err := h.URLService.RandomUrl(8)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate a random URL"})
		return
	}

	shortURL, err := h.URLService.SaveRandomUrl(req.OriginalURL, randomURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shorten URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Create successful": config.AppConfig.URL.SecretKey + shortURL})

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

	if !strings.HasPrefix(req.CustomShortURL, "http://") || !strings.HasPrefix(req.CustomShortURL, "https://") {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Invalid Url"})
		return
	}

	shortURl, err := h.URLService.SaveUrl(req.OriginalURL, req.CustomShortURL, userId.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shorten URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Create successful": "http://localhost:8080/api/shorten/" + shortURl})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	shortenedURL := c.Param("shortenedURL")

	originalURL, err := h.URLService.GetOriginalURL(shortenedURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found Url"})
		return
	}

	err = h.URLService.IncreaseClick(shortenedURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to increase click count"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)

}

func (h *URLHandler) GetId(c *gin.Context) {

	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Fail to get user"})
		return
	}

	UrlList, err := h.URLService.GetById(userId.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Noice": "You don't create any Url"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": UrlList})
}

func (h *URLHandler) Ratelimt(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Request successful"})
}

func (h *URLHandler) GetAllCacheUrl(c *gin.Context) {
	urls, err := h.URLService.GetAllUrlFromCache()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't retrieve Urls"})
		return
	}

	c.JSON(http.StatusOK, urls)
}
