package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const envPath = "./.env"

// DatabaseConfig represents the configuration for the database.
type DatabaseConfig struct {
	PostgresDB       string `env:"POSTGRES_DB" env-required:"true"`
	PostgresUser     string `env:"POSTGRES_USER" env-required:"true"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DbHost           string `env:"DB_HOST" env-required:"true"`
	DbPort           int    `env:"DB_PORT" env-required:"true"`
	PostgresDSN      string `env:"-"`
}

// GRPCConfig represents the configuration for the gRPC server.
type GRPCConfig struct {
	GRPCPort int `env:"GRPC_PORT" env-required:"true"`
}

// Config represents the overall application configuration.
type Config struct {
	AppName  string `env:"APP_NAME" env-required:"true"`
	Database DatabaseConfig
	GRPC     GRPCConfig
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

	cfg.Database.PostgresDSN = fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Database.DbHost,
		cfg.Database.DbPort,
		cfg.Database.PostgresDB,
		cfg.Database.PostgresUser,
		cfg.Database.PostgresPassword,
	)

	return &cfg, nil
}
