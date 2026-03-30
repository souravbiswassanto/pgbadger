package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/souravbiswassanto/pgbadger/config"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(cfg *config.Config, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		cfg:    cfg,
		logger: logger,
	}
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("Invalid login request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Verify username
	if req.Username != h.cfg.Auth.Username {
		h.logger.Warnw("Login attempt with invalid username", "username", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(h.cfg.Auth.PasswordHash), []byte(req.Password)); err != nil {
		h.logger.Warnw("Login attempt with invalid password", "username", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := h.generateToken(req.Username)
	if err != nil {
		h.logger.Errorw("Failed to generate token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	h.logger.Infow("User logged in successfully", "username", req.Username)
	c.JSON(http.StatusOK, LoginResponse{Token: token})
}

// AuthMiddleware validates JWT tokens
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for login and health endpoints
		if c.Request.URL.Path == "/api/v1/login" || c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(h.cfg.Auth.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			h.logger.Warnw("Invalid token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if username, exists := claims["username"]; exists {
				c.Set("username", username)
			}
		}

		c.Next()
	}
}

func (h *AuthHandler) generateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(h.cfg.Auth.JWTExpiry).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.Auth.JWTSecret))
}
