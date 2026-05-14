package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"rides-api/internal/metrics"
	"rides-api/internal/ride"
)

type Server struct {
	store   *ride.Store
	logger  *slog.Logger
	metrics *metrics.Metrics
	mux     *http.ServeMux
}

func NewServer(store *ride.Store, logger *slog.Logger, metrics *metrics.Metrics) *Server {
	s := &Server{
		store:   store,
		logger:  logger,
		metrics: metrics,
		mux:     http.NewServeMux(),
	}

	s.routes()

	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.handleHealth)
	s.mux.HandleFunc("GET /readyz", s.handleReady)
	s.mux.Handle("GET /metrics", s.metrics.Handler())
	s.mux.HandleFunc("GET /api/rides", s.handleListRides)
	s.mux.HandleFunc("POST /api/rides", s.handleCreateRide)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	recorder := &statusRecorder{
		ResponseWriter: w,
		status:         http.StatusOK,
	}

	s.mux.ServeHTTP(recorder, r)
	duration := time.Since(start)

	s.metrics.ObserveHTTPRequest(r.Method, r.URL.Path, recorder.status, duration.Seconds())

	s.logger.Info("http request",
		"method", r.Method,
		"path", r.URL.Path,
		"status", recorder.status,
		"duration_ms", duration.Milliseconds(),
	)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (s *Server) handleListRides(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"rides": s.store.List()})
}

func (s *Server) handleCreateRide(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req ride.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := ride.ValidateCreateRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created := s.store.Create(req)
	s.metrics.ObserveRideCreated()
	s.logger.Info("ride created",
		"ride_id", created.ID,
		"rider_id", created.RiderID,
		"pickup", created.Pickup,
		"dropoff", created.Dropoff,
	)
	writeJSON(w, http.StatusCreated, created)
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
		slog.Error("failed to write JSON response", "error", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
