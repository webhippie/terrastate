package handler

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "terrastate"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      "request_count_total",
			Help:      "counter of http requests made",
		},
		[]string{"action", "state"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      "request_duration_milliseconds",
			Help:      "histogram of the time (in milliseconds) each request took",
			Buckets:   append([]float64{.001, .003}, prometheus.DefBuckets...),
		},
		[]string{"action", "state"},
	)
)

func init() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}

func handleMetrics(start time.Time, action, state string) {
	duration := time.Since(start).Seconds() * 1e3

	requestCounter.WithLabelValues(action, state).Inc()
	requestDuration.WithLabelValues(action, state).Observe(duration)
}
