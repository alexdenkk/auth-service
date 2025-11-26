package config

import (
	"os"
	"time"
)

// Service configuration
type Config struct {
	DBConfig   *DBConfig
	JwtConfig  *JwtConfig
	HttpConfig *HttpConfig
}

// Database configuration
type DBConfig struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

// JWT configuration
type JwtConfig struct {
	SignKey []byte
}

// HTTP configuration
type HttpConfig struct {
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Loading configuration from environment
func NewConfigFromEnv() *Config {
	return &Config{
		// Database
		DBConfig: &DBConfig{
			Name:     os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},

		// JWT
		JwtConfig: &JwtConfig{
			SignKey: []byte(os.Getenv("JWT_SIGN_KEY")),
		},

		// HTTP
		HttpConfig: &HttpConfig{
			Host:         os.Getenv("HOST"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}
