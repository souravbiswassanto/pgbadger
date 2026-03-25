package logger

import (
	"go.uber.org/zap"
)

// New returns a SugaredLogger pointer configured with the provided level (info/debug)
func New(level string) *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	if level == "debug" {
		cfg = zap.NewDevelopmentConfig()
	}

	l, err := cfg.Build()
	if err != nil {
		// fallback to a no-op logger
		nol, _ := zap.NewProduction()
		return nol.Sugar()
	}

	return l.Sugar()
}
