package config

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
)

type Config struct {
	AppEnv        string
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
		AppEnv:        getEnv("APP_ENV", "development"),
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

// weakSecrets — плейсхолдеры/дефолты, которыми НЕЛЬЗЯ подписывать JWT в проде:
// они известны публично (лежат в репозитории), поэтому любой может подделать
// токен с произвольным user_id и войти кем угодно, включая superuser.
var weakSecrets = map[string]bool{
	"":                                 true,
	"dev-secret-key":                   true,
	"change-me-in-production-please":    true,
	"CHANGE_ME_RANDOM_STRING_64_CHARS": true,
}

func (c *Config) isDev() bool {
	switch c.AppEnv {
	case "development", "dev", "local", "test":
		return true
	}
	return false
}

// Validate падает на старте (fail-fast), если в production-окружении
// JWT_SECRET пустой или равен известному дефолту. В dev-окружении это лишь
// громкое предупреждение, чтобы локальная разработка работала из коробки.
func (c *Config) Validate() error {
	if weakSecrets[c.JWTSecret] {
		if !c.isDev() {
			return fmt.Errorf("JWT_SECRET пустой или равен известному дефолту — задайте сильный секрет (APP_ENV=%q)", c.AppEnv)
		}
		log.Printf("WARNING: JWT_SECRET слабый/дефолтный — допустимо только в dev (APP_ENV=%q)", c.AppEnv)
	}
	return nil
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
