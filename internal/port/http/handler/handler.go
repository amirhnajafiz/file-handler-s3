package handler

import (
	"net/http"

	"github.com/amirhnajafiz/hls/internal/port/s3"
	"github.com/amirhnajafiz/hls/internal/telemetry/metric"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	S3     s3.Handler
	Trace  trace.Tracer
	Metric metric.Metrics
}

// Home will return the home page
func (h Handler) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Metric.Requests.Add(1)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		http.ServeFile(w, r, "./templates/index.html")
	}
}

// Files will return the files page
func (h Handler) Files() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Metric.Requests.Add(1)

		w.Header().Set("Access-Control-Allow-Origin", "*")

		http.ServeFile(w, r, "./templates/files.html")
	}
}
