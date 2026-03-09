package config

import (
	"os"
	"strings"
)

type Config struct {
	Environment    string
	Port           string
	JWTSecret      string
	DatabaseURL    string
	NATSURL        string
	AWSRegion      string
	AWSAccountID   string
	SlackURL       string
	AllowedOrigins []string
}

func Load() Config {
	return Config{
		Environment:    getEnv("APP_ENV", "local"),
		Port:           getEnv("API_PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/cost_guard?sslmode=disable"),
		NATSURL:        getEnv("NATS_URL", "nats://localhost:4222"),
		AWSRegion:      getEnv("AWS_REGION", "us-east-1"),
		AWSAccountID:   getEnv("AWS_COST_EXPLORER_ACCOUNT_ID", ""),
		SlackURL:       getEnv("SLACK_WEBHOOK_URL", ""),
		AllowedOrigins: getEnvList("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:5173"}),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvList(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}

	if len(items) == 0 {
		return fallback
	}

	return items
}
