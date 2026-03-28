package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName         string
	AppEnv          string
	AppPort         string
	AllowedOrigins  string
	RateLimitMax    int
	RateLimitWindow int
	DBHost          string
	DBPort          string
	DBName          string
	DBUser          string
	DBPassword      string
	JWTSecret       string
	JWTExpiresHours int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	expiresHours := 24
	if raw := os.Getenv("JWT_EXPIRES_HOURS"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid JWT_EXPIRES_HOURS: %w", err)
		}
		expiresHours = parsed
	}

	rateLimitMax := 120
	if raw := os.Getenv("RATE_LIMIT_MAX"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid RATE_LIMIT_MAX: %w", err)
		}
		rateLimitMax = parsed
	}

	rateLimitWindow := 60
	if raw := os.Getenv("RATE_LIMIT_WINDOW_SECONDS"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid RATE_LIMIT_WINDOW_SECONDS: %w", err)
		}
		rateLimitWindow = parsed
	}

	cfg := &Config{
		AppName:         getOrDefault("APP_NAME", "canhbuomxanh-api"),
		AppEnv:          getOrDefault("APP_ENV", "development"),
		AppPort:         getOrDefault("APP_PORT", "8080"),
		AllowedOrigins:  getOrDefault("ALLOWED_ORIGINS", "http://localhost:5500,http://127.0.0.1:5500"),
		RateLimitMax:    rateLimitMax,
		RateLimitWindow: rateLimitWindow,
		DBHost:          getOrDefault("DB_HOST", "127.0.0.1"),
		DBPort:          getOrDefault("DB_PORT", "3306"),
		DBName:          getOrDefault("DB_NAME", "canhbuomxanh"),
		DBUser:          getOrDefault("DB_USER", "root"),
		DBPassword:      getOrDefault("DB_PASSWORD", "root"),
		JWTSecret:       getOrDefault("JWT_SECRET", "change_me_to_long_random_value"),
		JWTExpiresHours: expiresHours,
	}

	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func getOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
