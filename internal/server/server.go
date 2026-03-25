package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/souravbiswassanto/pgbadger/config"
	"github.com/souravbiswassanto/pgbadger/internal/handler"
	"github.com/souravbiswassanto/pgbadger/internal/middleware"
)

type Server struct {
	cfg  *config.Config
	lg   *zap.SugaredLogger
	http *http.Server
}

func New(cfg *config.Config, lg *zap.SugaredLogger) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ZapLogger(lg))

	// register routes
	handler.Register(r, lg)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return &Server{cfg: cfg, lg: lg, http: srv}
}

func (s *Server) Start() error {
	s.lg.Infof("starting server on %s", s.http.Addr)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.lg.Info("shutting down server")
	return s.http.Shutdown(ctx)
}
