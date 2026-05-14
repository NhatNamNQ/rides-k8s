package config

import "os"

type Config struct {
	Port        string
	LogLevel    string
	DatabaseURL string
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
