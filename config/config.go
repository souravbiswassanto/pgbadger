package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`

	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`

	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// Load loads configuration from viper (env, file)
func Load() *Config {
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("read_timeout", "5s")
	viper.SetDefault("write_timeout", "10s")

	var cfg Config
	// support durations as strings
	viper.RegisterAlias("read_timeout", "read_timeout")
	viper.RegisterAlias("write_timeout", "write_timeout")

	_ = viper.Unmarshal(&cfg)

	// convert durations if provided as strings
	if rt := viper.GetString("read_timeout"); rt != "" {
		if d, err := time.ParseDuration(rt); err == nil {
			cfg.ReadTimeout = d
		}
	}
	if wt := viper.GetString("write_timeout"); wt != "" {
		if d, err := time.ParseDuration(wt); err == nil {
			cfg.WriteTimeout = d
		}
	}

	return &cfg
}
