package main

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rides-simulator/internal/config"
	"rides-simulator/internal/sim"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel(),
	}))

	simulator := sim.New(cfg, logger, rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()+1))))

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      simulator.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("simulator loop starting",
			"api_base_url", cfg.APIBaseURL,
			"interval", cfg.Interval.String(),
		)
		simulator.Run(ctx)
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("failed to shut down simulator server", "error", err)
		}
	}()

	logger.Info("rides simulator listening", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("simulator server failed", "error", err)
		os.Exit(1)
	}
}
