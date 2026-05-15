package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"rides-api/internal/metrics"
	"rides-api/internal/ride"
)

func testServer() *Server {
	return NewServer(ride.NewMemoryStore(), slog.New(slog.NewTextHandler(io.Discard, nil)), metrics.New())
}

func TestHealthz(t *testing.T) {
	server := testServer()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestReadyz(t *testing.T) {
	server := testServer()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}

func TestReadyzWhenStoreIsUnavailable(t *testing.T) {
	server := NewServer(
		failingStore{},
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		metrics.New(),
	)
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)

	if res.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, res.Code)
	}
}

func TestCreateRide(t *testing.T) {
	server := testServer()
	body := bytes.NewBufferString(`{"rider_id":"rider-1","pickup":"zone-a","dropoff":"zone-b"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/rides", body)
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}

	var created ride.Ride
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if created.ID != "ride-1" {
		t.Fatalf("expected ride ID ride-1, got %q", created.ID)
	}
	if created.Status != "requested" {
		t.Fatalf("expected status requested, got %q", created.Status)
	}
}

func TestCreateRideValidation(t *testing.T) {
	server := testServer()
	body := bytes.NewBufferString(`{"rider_id":"","pickup":"zone-a","dropoff":"zone-b"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/rides", body)
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestListRides(t *testing.T) {
	server := testServer()

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/api/rides",
		bytes.NewBufferString(`{"rider_id":"rider-1","pickup":"zone-a","dropoff":"zone-b"}`),
	)
	createRes := httptest.NewRecorder()
	server.ServeHTTP(createRes, createReq)

	listReq := httptest.NewRequest(http.MethodGet, "/api/rides", nil)
	listRes := httptest.NewRecorder()
	server.ServeHTTP(listRes, listReq)

	if listRes.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listRes.Code)
	}

	var response struct {
		Rides []ride.Ride `json:"rides"`
	}
	if err := json.NewDecoder(listRes.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(response.Rides) != 1 {
		t.Fatalf("expected 1 ride, got %d", len(response.Rides))
	}
}

func TestMetricsEndpoint(t *testing.T) {
	server := testServer()

	createReq := httptest.NewRequest(
		http.MethodPost,
		"/api/rides",
		bytes.NewBufferString(`{"rider_id":"rider-1","pickup":"zone-a","dropoff":"zone-b"}`),
	)
	createRes := httptest.NewRecorder()
	server.ServeHTTP(createRes, createReq)

	metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsRes := httptest.NewRecorder()
	server.ServeHTTP(metricsRes, metricsReq)

	if metricsRes.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, metricsRes.Code)
	}

	body := metricsRes.Body.String()
	for _, expected := range []string{
		"rides_created_total 1",
		"rides_active 1",
		`http_requests_total{method="POST",path="/api/rides",status="201"} 1`,
	} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected metrics response to contain %q, got:\n%s", expected, body)
		}
	}
}

type failingStore struct{}

func (failingStore) Create(ctx context.Context, req ride.CreateRequest) (ride.Ride, error) {
	return ride.Ride{}, errors.New("store unavailable")
}

func (failingStore) List(ctx context.Context) ([]ride.Ride, error) {
	return nil, errors.New("store unavailable")
}

func (failingStore) Ping(ctx context.Context) error {
	return errors.New("store unavailable")
}
