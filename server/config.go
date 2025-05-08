// Package server implements the server for the terminal.
package server

import "time"

// Config holds the server configuration
type Config struct {
	// Server configuration
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`

	// JWT configuration
	JWTSecret            string        `mapstructure:"jwt_secret"`
	JWTExpirationTime    time.Duration `mapstructure:"jwt_expiration_time"`
	JWTRefreshExpiration time.Duration `mapstructure:"jwt_refresh_expiration"`

	// Database configuration
	DatabaseURL string `mapstructure:"database_url"`

	// CORS configuration
	CORSAllowOrigins []string `mapstructure:"cors_allow_origins"`

	// API rate limiting
	RateLimitEnabled bool          `mapstructure:"rate_limit_enabled"`
	RateLimitValue   int           `mapstructure:"rate_limit_value"`
	RateLimitWindow  time.Duration `mapstructure:"rate_limit_window"`
}

// DefaultConfig returns a config with default values
func DefaultConfig() *Config {
	return &Config{
		Host:                 "127.0.0.1",
		Port:                 8080,
		ReadTimeout:          5 * time.Second,
		WriteTimeout:         10 * time.Second,
		ShutdownTimeout:      30 * time.Second,
		JWTSecret:            "supersecretkey-change-me-in-production",
		JWTExpirationTime:    24 * time.Hour,
		JWTRefreshExpiration: 7 * 24 * time.Hour,
		CORSAllowOrigins:     []string{"*"},
		RateLimitEnabled:     true,
		RateLimitValue:       100,         // 100 requests
		RateLimitWindow:      time.Minute, // per minute
	}
}
