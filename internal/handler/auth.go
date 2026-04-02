package handler

import (
	"net/http"
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

	// Double Submit Cookie CSRF protection: require X-CSRF-Token header equals csrf_token cookie
	csrfHeader := c.GetHeader("X-CSRF-Token")
	csrfCookie, err := c.Cookie("csrf_token")
	if err != nil || csrfHeader == "" || csrfCookie == "" || csrfHeader != csrfCookie {
		h.logger.Warnw("CSRF token mismatch or missing", "header", csrfHeader, "cookie", csrfCookie)
		c.JSON(http.StatusForbidden, gin.H{"error": "CSRF validation failed"})
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

	// Set JWT in HttpOnly cookie (max age 1h)
	maxAge := int((time.Hour).Seconds())
	secure := h.cfg.Security.EnableTLS
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
		Expires:  time.Now().Add(time.Hour),
	}
	http.SetCookie(c.Writer, cookie)

	// TODO: maybe set different csrfVal
	// Also set a csrf cookie (double-submit). This cookie must NOT be HttpOnly so client JS can read it.
	csrfVal := csrfCookie // reuse the existing csrf cookie value set when serving login page
	csrf := &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfVal,
		Path:     "/",
		HttpOnly: false,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
		Expires:  time.Now().Add(time.Hour),
	}
	http.SetCookie(c.Writer, csrf)

	h.logger.Infow("User logged in successfully", "username", req.Username)
	// return success (no token in body)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	// TODO: Maybe redirect
}

// AuthMiddleware validates JWT tokens
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for login page and health endpoints
		if c.Request.URL.Path == "/api/v1/login" || c.Request.URL.Path == "/health" || c.Request.URL.Path == "/login" {
			c.Next()
			return
		}

		// If insecure mode is enabled, bypass authentication
		if h.cfg.Security.Insecure {
			c.Set("username", "insecure")
			c.Next()
			return
		}

		// Read token from cookie
		tokenStr, err := c.Cookie("auth_token")
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

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

		// CSRF check for state-changing requests (double-submit cookie)
		method := c.Request.Method
		if method != http.MethodGet && method != http.MethodHead && method != http.MethodOptions {
			csrfHeader := c.GetHeader("X-CSRF-Token")
			csrfCookie, _ := c.Cookie("csrf_token")
			if csrfHeader == "" || csrfCookie == "" || csrfHeader != csrfCookie {
				h.logger.Warnw("CSRF token mismatch or missing", "path", c.Request.URL.Path)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "CSRF validation failed"})
				return
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
