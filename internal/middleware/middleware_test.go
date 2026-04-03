package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeaders(true))
	r.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	// spot check a few headers
	require.NotEmpty(t, w.Header().Get("Content-Security-Policy"))
	require.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	require.NotEmpty(t, w.Header().Get("Strict-Transport-Security"))
}
