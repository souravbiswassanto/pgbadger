package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

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
	// security headers (HSTS added when TLS is enabled)
	r.Use(middleware.SecurityHeaders(cfg.Security.EnableTLS))
	r.Use(middleware.ZapLogger(lg))

	// serve static web assets (login.js, etc.) under /static
	r.Static("/static", "./web")

	// register routes
	handler.Register(r, cfg, lg)

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
	if s.cfg != nil && s.cfg.Security.EnableTLS {
		if s.cfg.Security.CertFile == "" || s.cfg.Security.KeyFile == "" {
			s.lg.Fatal("TLS enabled but cert_file or key_file is not provided in config")
		}
		s.lg.Infof("starting TLS server with cert=%s key=%s", s.cfg.Security.CertFile, s.cfg.Security.KeyFile)

		// Prepare TLS config and optionally append CA file into cert pool
		var tlsCfg tls.Config

		// Load server certificate
		cert, err := tls.LoadX509KeyPair(s.cfg.Security.CertFile, s.cfg.Security.KeyFile)
		if err != nil {
			s.lg.Fatalf("failed to load server cert/key: %v", err)
		}
		tlsCfg.Certificates = []tls.Certificate{cert}

		// If CA file is provided, append it to system pool and use as RootCAs
		if s.cfg.Security.CAFile != "" {
			caPEM, err := os.ReadFile(s.cfg.Security.CAFile)
			if err != nil {
				s.lg.Fatalf("failed to read CA file: %v", err)
			}

			pool, err := x509.SystemCertPool()
			if err != nil || pool == nil {
				pool = x509.NewCertPool()
			}

			if ok := pool.AppendCertsFromPEM(caPEM); !ok {
				s.lg.Warn("no certs appended from CA file; check PEM format")
			}

			tlsCfg.RootCAs = pool
		}

		s.http.TLSConfig = &tlsCfg

		// Use ListenAndServeTLS with empty cert paths because certificates are set in TLSConfig
		return s.http.ListenAndServeTLS("", "")
	}

	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.lg.Info("shutting down server")
	return s.http.Shutdown(ctx)
}
