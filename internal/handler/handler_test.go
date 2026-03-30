package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/souravbiswassanto/pgbadger/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Create minimal config and logger for testing
	cfg := &config.Config{}
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	Register(r, cfg, sugar)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}
