package sim

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"rides-simulator/internal/config"
)

type fixedRandom struct{}

func (fixedRandom) IntN(n int) int {
	return 0
}

func TestCreateRide(t *testing.T) {
	t.Parallel()

	requests := 0
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/rides" {
			t.Fatalf("expected /api/rides, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer api.Close()

	simulator := New(config.Config{
		APIBaseURL: api.URL,
		Interval:   time.Millisecond,
	}, slog.New(slog.NewTextHandler(io.Discard, nil)), fixedRandom{})

	if err := simulator.createRide(context.Background()); err != nil {
		t.Fatalf("expected createRide to succeed, got error: %v", err)
	}

	if requests != 1 {
		t.Fatalf("expected exactly one request, got %d", requests)
	}
}
