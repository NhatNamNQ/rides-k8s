package sim

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"rides-simulator/internal/config"
)

var zones = []string{"zone-a", "zone-b", "zone-c", "zone-d", "zone-e"}

type randomSource interface {
	IntN(n int) int
}

type Simulator struct {
	cfg        config.Config
	logger     *slog.Logger
	client     *http.Client
	metrics    *Metrics
	random     randomSource
	httpServer http.Handler
}

type createRideRequest struct {
	RiderID string `json:"rider_id"`
	Pickup  string `json:"pickup"`
	Dropoff string `json:"dropoff"`
}

func New(cfg config.Config, logger *slog.Logger, random randomSource) *Simulator {
	s := &Simulator{
		cfg:    cfg,
		logger: logger,
		client: &http.Client{Timeout: 10 * time.Second},
		metrics: NewMetrics(),
		random: random,
	}

	s.httpServer = s.routes()
	return s
}

func (s *Simulator) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		s.httpServer.ServeHTTP(recorder, r)
		duration := time.Since(start)

		s.metrics.ObserveHTTPRequest(r.Method, r.URL.Path, recorder.status, duration.Seconds())
	})
}

func (s *Simulator) Run(ctx context.Context) {
	interval := time.NewTicker(s.cfg.Interval)
	defer interval.Stop()

	for {
		loopStart := time.Now()

		if err := s.createRide(ctx); err != nil {
			s.logger.Error("simulator failed to create ride", "error", err)
			s.metrics.ObserveEvent("create_ride", "error")
		} else {
			s.metrics.ObserveEvent("create_ride", "success")
		}

		s.metrics.ObserveLoopDuration(time.Since(loopStart).Seconds())

		select {
		case <-ctx.Done():
			s.logger.Info("simulator loop stopped")
			return
		case <-interval.C:
		}
	}
}

func (s *Simulator) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /metrics", s.metrics.Handler())
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	return mux
}

func (s *Simulator) createRide(ctx context.Context) error {
	payload := createRideRequest{
		RiderID: fmt.Sprintf("sim-rider-%d", time.Now().UnixNano()),
		Pickup:  pickZone(s.random),
		Dropoff: pickZone(s.random),
	}

	if payload.Pickup == payload.Dropoff {
		payload.Dropoff = zones[(indexOf(payload.Dropoff)+1)%len(zones)]
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.APIBaseURL+"/api/rides", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return fmt.Errorf("post ride: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status: %d", response.StatusCode)
	}

	s.logger.Info("simulator created ride",
		"rider_id", payload.RiderID,
		"pickup", payload.Pickup,
		"dropoff", payload.Dropoff,
	)

	return nil
}

func pickZone(random randomSource) string {
	return zones[random.IntN(len(zones))]
}

func indexOf(zone string) int {
	for index, candidate := range zones {
		if candidate == zone {
			return index
		}
	}
	return 0
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to write simulator JSON response", "error", err)
	}
}
