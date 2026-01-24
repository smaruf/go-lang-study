package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port   string
	Host   string
	DBPath string
	Env    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		Port:   getEnv("PORT", "8080"),
		Host:   getEnv("HOST", "localhost"),
		DBPath: getEnv("DB_PATH", "./data/tasks.db"),
		Env:    getEnv("ENV", "development"),
	}

	return config, nil
}

// getEnv retrieves an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Address returns the full server address
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
