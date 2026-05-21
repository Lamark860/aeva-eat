package config

import (
	"net"
	"net/url"
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
		GeosuggestKey: getEnv("GEOSUGGEST_KEY", ""),
	}
}

func (c *Config) DatabaseURL() string {
	// url.URL экранирует логин/пароль — иначе спецсимволы в пароле
	// (например `(`, `@`, `/`) ломают парсинг DSN.
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Host:   net.JoinHostPort(c.DBHost, c.DBPort),
		Path:   c.DBName,
	}
	q := url.Values{}
	q.Set("sslmode", c.DBSSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
