package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "divelog_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "divelog_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// wrap the response writer so we can capture the status code
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start).Seconds()
		path := r.URL.Path
		method := r.Method
		status := strconv.Itoa(wrapped.status)

		requestCount.WithLabelValues(method, path, status).Inc()
		requestDuration.WithLabelValues(method, path).Observe(duration)
	})
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}
