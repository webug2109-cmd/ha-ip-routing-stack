package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response latency (seconds).",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func lookupHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := http.StatusOK
	path := "/lookup"

	defer func() {
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(path).Observe(duration)
		httpRequestsTotal.WithLabelValues(path, fmt.Sprintf("%d", status)).Inc()
	}()

	ip := strings.TrimPrefix(r.URL.Path, "/lookup/")
	if ip == "" {
		status = http.StatusBadRequest
		http.Error(w, "IP address is required", status)
		return
	}

	// Make a real call to IPinfo free API
	resp, err := http.Get(fmt.Sprintf("https://ipinfo.io/%s/geo", ip))
	if err != nil {
		status = http.StatusInternalServerError
		http.Error(w, "Failed to query IPinfo", status)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func main() {
	http.HandleFunc("/lookup/", lookupHandler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Starting mock API server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
