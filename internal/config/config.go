package config

import (
	"time"

	"github.com/wb-go/wbf/config"
)

type Config struct {
	HTTPPort        string
	BaseURL         string
	PostgresDSN     string
	DBMaxConns      int
	DBConnAttempts  int
	DBRetryDelay    time.Duration
	DBMaxRetryDelay time.Duration
	GinMode         string
}

func Load() *Config {
	cfg := config.New()
	_ = cfg.LoadEnvFiles("configs/config.env")
	cfg.EnableEnv("")

	cfg.SetDefault("HTTP_PORT", "8080")
	cfg.SetDefault("BASE_URL", "http://localhost:8080")
	cfg.SetDefault("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable")
	cfg.SetDefault("DB_MAX_CONNS", 25)
	cfg.SetDefault("DB_CONN_ATTEMPTS", 10)
	cfg.SetDefault("DB_RETRY_DELAY", 100*time.Millisecond)
	cfg.SetDefault("DB_MAX_RETRY_DELAY", 5*time.Second)
	cfg.SetDefault("GIN_MODE", "debug")

	return &Config{
		HTTPPort:        cfg.GetString("HTTP_PORT"),
		BaseURL:         cfg.GetString("BASE_URL"),
		PostgresDSN:     cfg.GetString("POSTGRES_DSN"),
		DBMaxConns:      cfg.GetInt("DB_MAX_CONNS"),
		DBConnAttempts:  cfg.GetInt("DB_CONN_ATTEMPTS"),
		DBRetryDelay:    cfg.GetDuration("DB_RETRY_DELAY"),
		DBMaxRetryDelay: cfg.GetDuration("DB_MAX_RETRY_DELAY"),
		GinMode:         cfg.GetString("GIN_MODE"),
	}
}
