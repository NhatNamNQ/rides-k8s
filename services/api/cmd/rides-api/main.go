package main

import (
	"log/slog"
	"net/http"
	"os"

	"rides-api/internal/config"
	"rides-api/internal/httpapi"
	"rides-api/internal/ride"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))

	server := httpapi.NewServer(ride.NewStore(), logger)
	addr := ":" + cfg.Port

	logger.Info("rides API starting",
		"addr", addr,
		"database_configured", cfg.DatabaseURL != "",
	)

	if err := http.ListenAndServe(addr, server); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func parseLogLevel(level string) slog.Level {
	switch level {
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
