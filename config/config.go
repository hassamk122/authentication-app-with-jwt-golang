package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	DatabaseUrl string
	Environment string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("Error loading env %v", err)
	}

	return &Config{
		ServerPort:  getEnvOrDefault("SERVER_PORT", "8080"),
		DatabaseUrl: getEnvOrDefault("DB_URL", "postgres"),
		Environment: getEnvOrDefault("ENVIRONMENT", "dev"),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
