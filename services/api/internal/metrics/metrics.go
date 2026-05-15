package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	registry            *prometheus.Registry
	httpRequestsTotal   *prometheus.CounterVec // CounterVec is used to track the total number of requests, with labels for method, path, and status
	httpRequestDuration *prometheus.HistogramVec // Histogram is used to track the speed of requests
	ridesCreatedTotal   prometheus.Counter // just only increment when a ride is created
	ridesActive         prometheus.Gauge // increment when a ride is created, decrement when a ride is completed
}

func New() *Metrics {
	m := &Metrics{
		registry: prometheus.NewRegistry(),
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests handled by the API.",
			},
			[]string{"method", "path", "status"},
		),
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
		ridesCreatedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "rides_created_total",
				Help: "Total number of rides created.",
			},
		),
		ridesActive: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "rides_active",
				Help: "Current number of active rides.",
			},
		),
	}

	m.registry.MustRegister(
		m.httpRequestsTotal,
		m.httpRequestDuration,
		m.ridesCreatedTotal,
		m.ridesActive,
	)

	return m
}

func (m *Metrics) Handler() http.Handler { // This function make registry to HTTP handler, so that we can expose it endpoint
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

func (m *Metrics) ObserveHTTPRequest(method, path string, status int, durationSeconds float64) {
	labels := prometheus.Labels{
		"method": method,
		"path":   path,
		"status": strconv.Itoa(status),
	}

	m.httpRequestsTotal.With(labels).Inc()
	m.httpRequestDuration.With(labels).Observe(durationSeconds)
}

func (m *Metrics) ObserveRideCreated() {
	m.ridesCreatedTotal.Inc()
	m.ridesActive.Inc()
}
