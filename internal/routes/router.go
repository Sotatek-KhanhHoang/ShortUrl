package routes

import (
	"ShortUrl/internal/handlers"
	"ShortUrl/internal/services"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(authService *services.AuthService, urlService *services.URLService, redisClient *redis.Client) gin.Engine {
	r := gin.Default()

	authHandler := handlers.AuthHandler{AuthService: authService}
	urlHandler := handlers.URLHandler{URLService: urlService}
	rateLimiter := NewRateLimiter(redisClient, 10, time.Minute)

	auth := r.Group("/auth")
	auth.Use(rateLimiter.LimitMiddleWare())
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.GET("/:shortenedURL", urlHandler.Redirect)
		auth.POST("/shorten", urlHandler.CreateRandomUrl)
		auth.GET("/cache", urlHandler.GetAllCacheUrl)
	}

	api := r.Group("/api")
	api.Use(authMiddleware())
	api.Use(rateLimiter.LimitMiddleWare())
	{
		api.POST("/shorten", urlHandler.CreateShortenUrl)
		api.GET("/shorten", urlHandler.GetId)
		api.GET("/limit", urlHandler.Ratelimt)
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

type RateLimit struct {
	redisClient *redis.Client
	rate        int
	window      time.Duration
}

func NewRateLimiter(redisClient *redis.Client, rate int, window time.Duration) *RateLimit {
	return &RateLimit{
		redisClient: redisClient,
		rate:        rate,
		window:      window,
	}
}

func (rl *RateLimit) LimitMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		ipAddress := c.ClientIP()
		if ipAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "IP address is required"})
			c.Abort()
			return
		}

		key := "rate_limit:" + ipAddress

		pipe := rl.redisClient.TxPipeline()
		increase := pipe.Incr(context.Background(), key)
		pipe.Expire(context.Background(), key, rl.window)
		_, err := pipe.Exec(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to access Redis"})
			c.Abort()
			return
		}

		if increase.Val() > int64(rl.rate) {
			c.JSON(http.StatusTooManyRequests, gin.H{"erro": "Too many request"})
			c.Abort()
			return
		}

		c.Next()
	}
}
