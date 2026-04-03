package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/souravbiswassanto/pgbadger/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// prepare config with known username/password
	cfg := &config.Config{}
	cfg.Auth.Username = "testuser"
	// bcrypt hash for password "secret"
	h, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	require.NoError(t, err)
	cfg.Auth.PasswordHash = string(h)

	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	Register(r, cfg, sugar)

	// 1) GET /login should set csrf_token cookie
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	// cookie should be set
	cookies := w.Result().Cookies()
	var csrfVal string
	for _, c := range cookies {
		if c.Name == "csrf_token" {
			csrfVal = c.Value
		}
	}
	require.NotEmpty(t, csrfVal)

	// 2) POST /api/v1/login without CSRF header should be 403
	creds := map[string]string{"username": "testuser", "password": "secret"}
	body, _ := json.Marshal(creds)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// set cookie from previous response
	req.Header.Set("Cookie", "csrf_token="+csrfVal)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)

	// 3) POST with matching CSRF header should succeed and set auth_token cookie
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CSRF-Token", csrfVal)
	req.Header.Set("Cookie", "csrf_token="+csrfVal)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// auth_token cookie should be set
	cookies = w.Result().Cookies()
	var authToken string
	for _, c := range cookies {
		if c.Name == "auth_token" {
			authToken = c.Value
		}
	}
	require.NotEmpty(t, authToken)
}
