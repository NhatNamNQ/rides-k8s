package sim

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	registry               *prometheus.Registry
	httpRequestsTotal      *prometheus.CounterVec
	httpRequestDuration    *prometheus.HistogramVec
	simulatorEventsTotal   *prometheus.CounterVec
	simulatorLoopDuration  prometheus.Histogram
}

func NewMetrics() *Metrics {
	m := &Metrics{
		registry: prometheus.NewRegistry(),
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests handled by the simulator.",
			},
			[]string{"method", "path", "status"},
		),
		httpRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds for the simulator service.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
		simulatorEventsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "simulator_events_total",
				Help: "Total number of simulator events by type and result.",
			},
			[]string{"event", "result"},
		),
		simulatorLoopDuration: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "simulator_loop_duration_seconds",
				Help:    "Duration of one simulator loop iteration in seconds.",
				Buckets: prometheus.DefBuckets,
			},
		),
	}

	m.registry.MustRegister(
		m.httpRequestsTotal,
		m.httpRequestDuration,
		m.simulatorEventsTotal,
		m.simulatorLoopDuration,
	)

	return m
}

func (m *Metrics) Handler() http.Handler {
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

func (m *Metrics) ObserveEvent(event, result string) {
	m.simulatorEventsTotal.WithLabelValues(event, result).Inc()
}

func (m *Metrics) ObserveLoopDuration(durationSeconds float64) {
	m.simulatorLoopDuration.Observe(durationSeconds)
}
