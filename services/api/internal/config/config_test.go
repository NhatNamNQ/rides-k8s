package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("DATABASE_URL", "")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.Port)
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("expected default log level info, got %q", cfg.LogLevel)
	}
	if cfg.DatabaseURL != "" {
		t.Fatalf("expected empty database URL, got %q", cfg.DatabaseURL)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("DATABASE_URL", "postgres://example")

	cfg := Load()

	if cfg.Port != "9090" {
		t.Fatalf("expected port 9090, got %q", cfg.Port)
	}
	if cfg.LogLevel != "debug" {
		t.Fatalf("expected log level debug, got %q", cfg.LogLevel)
	}
	if cfg.DatabaseURL != "postgres://example" {
		t.Fatalf("expected configured database URL, got %q", cfg.DatabaseURL)
	}
}
