package middleware

/*import (
	"ShortUrl/internal/handlers"
	"ShortUrl/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRepuired(authService *services.AuthService) gin.HandlerFunc {
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
}*/
