package routes

import (
	"ShortUrl/internal/handlers"
	"ShortUrl/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authService *services.AuthService, urlService *services.URLService) gin.Engine {
	r := gin.Default()

	authHandler := handlers.AuthHandler{AuthService: authService}
	urlHandler := handlers.URLHandler{URLService: urlService}

	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	api := r.Group("/api")
	api.Use(authMiddleware(authService))
	{
		api.POST("/shorten", urlHandler.CreateShortenUrl)
	}
	return *r
}

func authMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		userID, err := handlers.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("userID", userID.String())
		c.Next()
	}
}
