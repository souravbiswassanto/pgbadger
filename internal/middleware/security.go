package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersLegacy returns middleware that sets common security headers.
// Deprecated: SecurityHeaders is implemented in middleware.go. This legacy
// function is kept to avoid build errors while migrating.
func SecurityHeadersLegacy(enableHSTS bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy: allow scripts only from same origin
		// Adjust as needed for your app (for example to allow analytics)
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'self';")

		// Clickjacking protection
		c.Writer.Header().Set("X-Frame-Options", "DENY")

		// MIME sniffing
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// Referrer policy
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions policy - restrict features
		c.Writer.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		if enableHSTS {
			// HSTS for 1 year
			c.Writer.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=%d; includeSubDomains; preload", 31536000))
		}

		// Proceed
		c.Next()
	}
}
