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
		auth.GET("/:shortenedURL", urlHandler.Redirect)
	}

	api := r.Group("/api")
	api.Use(authMiddleware())
	{
		api.POST("/shorten", urlHandler.CreateShortenUrl)
	}
	return *r
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not found"})
			c.Abort()
			return
		}

		// Nếu sử dụng schema "Bearer", tách token ra
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		userID, err := handlers.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorizeed"})
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
