package middleware

import (
	"strconv"
	"time"

	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

// MetricsCollector interface for collecting metrics
type MetricsCollector interface {
	IncrementRequestCount(method, path string, statusCode int)
	RecordRequestDuration(method, path string, duration time.Duration)
	IncrementErrorCount(method, path string, errorType string)
}

// SimpleMetricsCollector is a simple in-memory metrics collector
// Can be replaced with Prometheus or other metrics backends
type SimpleMetricsCollector struct {
	requestCounts    map[string]int64
	requestDurations map[string][]time.Duration
	errorCounts      map[string]int64
}

// NewSimpleMetricsCollector creates a new simple metrics collector
func NewSimpleMetricsCollector() *SimpleMetricsCollector {
	return &SimpleMetricsCollector{
		requestCounts:    make(map[string]int64),
		requestDurations: make(map[string][]time.Duration),
		errorCounts:      make(map[string]int64),
	}
}

// IncrementRequestCount increments request count
func (m *SimpleMetricsCollector) IncrementRequestCount(method, path string, statusCode int) {
	key := method + ":" + path + ":" + strconv.Itoa(statusCode)
	m.requestCounts[key]++
}

// RecordRequestDuration records request duration
func (m *SimpleMetricsCollector) RecordRequestDuration(method, path string, duration time.Duration) {
	key := method + ":" + path
	m.requestDurations[key] = append(m.requestDurations[key], duration)
	// Keep only last 100 durations per endpoint
	if len(m.requestDurations[key]) > 100 {
		m.requestDurations[key] = m.requestDurations[key][1:]
	}
}

// IncrementErrorCount increments error count
func (m *SimpleMetricsCollector) IncrementErrorCount(method, path string, errorType string) {
	key := method + ":" + path + ":" + errorType
	m.errorCounts[key]++
}

// MetricsMiddleware creates a metrics collection middleware
func MetricsMiddleware(collector MetricsCollector) gin.HandlerFunc {
	// Use simple collector if none provided
	if collector == nil {
		collector = NewSimpleMetricsCollector()
		logger.Component("auth.middleware.metrics").
			Info().
			Msg("using simple in-memory metrics collector")
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Collect metrics
		collector.IncrementRequestCount(method, path, statusCode)
		collector.RecordRequestDuration(method, path, duration)

		// Track errors (4xx, 5xx)
		if statusCode >= 400 {
			errorType := "client_error"
			if statusCode >= 500 {
				errorType = "server_error"
			}
			collector.IncrementErrorCount(method, path, errorType)
		}

		// Log slow requests
		if duration > 1*time.Second {
			logger.Component("auth.middleware.metrics").
				Warn().
				Str("method", method).
				Str("path", path).
				Int("status", statusCode).
				Dur("duration", duration).
				Msg("slow request detected")
		}
	}
}

