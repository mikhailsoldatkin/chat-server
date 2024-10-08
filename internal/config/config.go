package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const envPath = "./.env"

// DB represents the configuration for the database.
type DB struct {
	PostgresDB       string `env:"POSTGRES_DB" env-required:"true"`
	PostgresUser     string `env:"POSTGRES_USER" env-required:"true"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Host             string `env:"DB_HOST" env-required:"true"`
	Port             int    `env:"DB_PORT" env-required:"true"`
	PostgresDSN      string `env:"-"`
}

// GRPC represents the configuration for the gRPC server.
type GRPC struct {
	Host    string `env:"GRPC_HOST" env-required:"true"`
	Port    int    `env:"GRPC_PORT" env-required:"true"`
	Address string `env:"-"`
}

// Auth represents the configuration for the authentication server.
type Auth struct {
	Host    string `env:"AUTH_HOST" env-required:"true"`
	Port    int    `env:"AUTH_PORT" env-required:"true"`
	Address string `env:"-"`
}

// Logger represents configuration for logger.
type Logger struct {
	Level      string `env:"LOG_LEVEL" env-required:"true"`
	Filename   string `env:"LOG_FILENAME" env-required:"true"`
	MaxSizeMB  int    `env:"LOG_MAX_SIZE_MB" env-required:"true"`
	MaxBackups int    `env:"LOG_MAX_BACKUPS" env-required:"true"`
	MaxAgeDays int    `env:"LOG_MAX_AGE_DAYS" env-required:"true"`
}

// Jaeger represents the configuration for the Jaeger server.
type Jaeger struct {
	Host    string `env:"JAEGER_HOST" env-required:"true"`
	Port    int    `env:"JAEGER_PORT" env-required:"true"`
	Address string `env:"-"`
}

// Config represents the overall application configuration.
type Config struct {
	DB     DB
	GRPC   GRPC
	Auth   Auth
	Logger Logger
	Jaeger Jaeger
}

// Load reads configuration from .env file.
func Load() (*Config, error) {
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil, fmt.Errorf(".env file does not exist in project's root")
	}

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("cannot read config from environment variables: %w", err)
	}

	cfg.DB.PostgresDSN = fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.PostgresDB,
		cfg.DB.PostgresUser,
		cfg.DB.PostgresPassword,
	)

	cfg.GRPC.Address = fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
	cfg.Auth.Address = fmt.Sprintf("%s:%d", cfg.Auth.Host, cfg.Auth.Port)
	cfg.Jaeger.Address = fmt.Sprintf("%s:%d", cfg.Jaeger.Host, cfg.Jaeger.Port)

	return &cfg, nil
}
