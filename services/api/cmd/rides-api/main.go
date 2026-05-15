package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"rides-api/internal/config"
	"rides-api/internal/httpapi"
	"rides-api/internal/metrics"
	"rides-api/internal/ride"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))

	store, cleanup, err := buildStore(cfg, logger)
	if err != nil {
		logger.Error("failed to initialize store", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	server := httpapi.NewServer(store, logger, metrics.New())
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

func buildStore(cfg config.Config, logger *slog.Logger) (ride.Repository, func(), error) {
	if cfg.DatabaseURL == "" {
		logger.Warn("DATABASE_URL is empty; using in-memory store")
		return ride.NewMemoryStore(), func() {}, nil
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, nil, fmt.Errorf("open postgres: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	store := ride.NewPostgresStore(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := store.Ping(ctx); err != nil {
		db.Close()
		return nil, nil, err
	}

	logger.Info("connected to postgres")

	return store, func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close postgres", "error", err)
		}
	}, nil
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
