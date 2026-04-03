package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapLogger returns a simple request logger using zap SugaredLogger
func ZapLogger(lg *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		if raw != "" {
			path = path + "?" + raw
		}

		lg.Infow("request",
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"latency", latency.String(),
			"client_ip", c.ClientIP(),
		)
	}
}

// SecurityHeaders sets a small set of recommended security headers including CSP
// SecurityHeaders returns middleware that sets common security headers.
// If enableHSTS is true, the caller should pass true to add Strict-Transport-Security.
func SecurityHeaders(enableHSTS bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy: allow scripts only from same origin
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none'; base-uri 'self';")

		// Clickjacking protection
		c.Header("X-Frame-Options", "DENY")

		// MIME sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Disable legacy XSS protection
		c.Header("X-XSS-Protection", "0")

		if enableHSTS {
			// HSTS for 1 year
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		c.Next()
	}
}
