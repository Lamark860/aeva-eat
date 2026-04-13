package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	APIPort       string
	GeosuggestKey string
}

func Load() *Config {
	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "aeva"),
		DBPassword:    getEnv("DB_PASSWORD", "aeva_secret"),
		DBName:        getEnv("DB_NAME", "aeva_eat"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		JWTSecret:     getEnv("JWT_SECRET", "dev-secret-key"),
		APIPort:       getEnv("API_PORT", "8085"),
		GeosuggestKey: getEnv("GEOSUGGEST_KEY", "9500e2ca-ad3b-4503-be48-8b732a773857"),
	}
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
