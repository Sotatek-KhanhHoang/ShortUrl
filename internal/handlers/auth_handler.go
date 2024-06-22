package handlers

import (
	"ShortUrl/internal/config"
	"ShortUrl/internal/services"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	//tạo biến để lưu dữ liệu request
	var loginReq LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid credential"))
	}

	fmt.Println(loginReq)
	//Gọi phương thức Login với dữ liệu của request
	user, err := h.AuthService.Login(loginReq.Username, loginReq.Password)
	//Xử lý lỗi
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentialss"})
		return
	}

	tokenS, err := generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie := http.Cookie{
		Name:     "user_id",
		Value:    tokenS,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(c.Writer, &cookie)

	c.Header("Authorization", "Bearer "+tokenS)

	c.JSON(http.StatusOK, gin.H{"message": tokenS})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.AuthService.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

func generateJWT(useriD uuid.UUID) (string, error) {

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: useriD,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.AppConfig.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ValidateToken(tokenStr string) (uuid.UUID, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.SecretKey), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("Invalid Tokens")
	}

	return claims.UserID, nil
}
