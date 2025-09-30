package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the configuration for the birb-client service
type Config struct {
	// BirbNestURL is the base URL of the birb-nest API
	BirbNestURL string

	// WriteInterval is the duration between writes
	WriteInterval time.Duration

	// LogLevel controls the logging verbosity
	LogLevel string
}

// Load loads the configuration from environment variables with sensible defaults
func Load() *Config {
	return &Config{
		BirbNestURL:   getEnv("BIRB_NEST_URL", "http://localhost:8080"),
		WriteInterval: getDurationEnv("WRITE_INTERVAL", 3*time.Second),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv retrieves a duration from environment (in seconds) or returns a default
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}