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

	Auth struct {
		JWTSecret    string        `mapstructure:"jwt_secret"`
		JWTExpiry    time.Duration `mapstructure:"jwt_expiry"`
		Username     string        `mapstructure:"username"`
		PasswordHash string        `mapstructure:"password_hash"`
	} `mapstructure:"auth"`

	Security struct {
		EnableTLS bool   `mapstructure:"enable_tls"`
		CertFile  string `mapstructure:"cert_file"`
		KeyFile   string `mapstructure:"key_file"`
		CAFile    string `mapstructure:"ca_file"`
		Insecure  bool   `mapstructure:"insecure"` // if true, bypass auth checks (testing only)
	} `mapstructure:"security"`

	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

/*
# Map structure
## Config file
server:
	host: 127.0.0.1
	port: 3000
log:
	level: debug
read_timeout: 10s
write_timeout: 30s

## env

SERVER_HOST=127.0.0.1
SERVER_PORT=3000
LOG_LEVEL=debug
READ_TIMEOUT=10s

viper.Unmarshal(&config)
*/

// Load loads configuration from viper (env, file)
func Load() *Config {
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 2385)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("auth.jwt_secret", "change-this-secret-in-production")
	viper.SetDefault("auth.jwt_expiry", "24h")
	viper.SetDefault("security.enable_tls", false)
	viper.SetDefault("security.cert_file", "")
	viper.SetDefault("security.key_file", "")
	viper.SetDefault("security.ca_file", "")
	viper.SetDefault("security.insecure", false)
	viper.SetDefault("auth.username", "admin")
	viper.SetDefault("auth.password_hash", "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi") // "password"
	viper.SetDefault("read_timeout", "5s")
	viper.SetDefault("write_timeout", "10s")

	var cfg Config
	//// support durations as strings
	//viper.RegisterAlias("read_timeout", "read_timeout")
	//viper.RegisterAlias("write_timeout", "write_timeout")

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

	// parse auth.jwt_expiry which may be provided as string like "24h"
	if je := viper.GetString("auth.jwt_expiry"); je != "" {
		if d, err := time.ParseDuration(je); err == nil {
			cfg.Auth.JWTExpiry = d
		}
	}

	return &cfg
}
