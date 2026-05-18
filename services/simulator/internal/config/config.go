package config

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

type Config struct {
	Port       string
	LogLevelRaw string
	APIBaseURL string
	Interval   time.Duration
}

func Load() Config {
	return Config{
		Port:        getEnv("PORT", "8081"),
		LogLevelRaw: getEnv("LOG_LEVEL", "info"),
		APIBaseURL:  strings.TrimRight(getEnv("API_BASE_URL", "http://localhost:8080"), "/"),
		Interval:    parseDuration(getEnv("SIMULATOR_INTERVAL", "5s"), 5*time.Second),
	}
}

func (c Config) LogLevel() slog.Level {
	switch strings.ToLower(c.LogLevelRaw) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func parseDuration(value string, fallback time.Duration) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return duration
}
