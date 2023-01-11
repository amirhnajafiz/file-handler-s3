package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amirhnajafiz/fhs/internal/port/http/handler"
	"github.com/amirhnajafiz/fhs/internal/port/s3"
	"github.com/amirhnajafiz/fhs/internal/telemetry/metric"
	"go.opentelemetry.io/otel/trace"
)

func New(cfg Config, t trace.Tracer, m metric.Metrics, sh s3.Handler) {
	// handler init
	h := handler.Handler{
		Trace:  t,
		Metric: m,
		S3:     sh,
	}

	// root
	http.Handle("/", h.Home())
	// files
	http.Handle("/files", h.Files())
	// add a handler for uploading a file
	http.Handle("/upload", h.UploadFile())
	// get all files method
	http.Handle("/all", h.GetAllFiles())
	// remove a file
	http.Handle("/del", h.RemoveFile())
	// download file
	http.Handle("/get", h.DownloadFile())

	log.Printf("Starting server on %v\n", cfg.Port)

	// serve and log errors
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), nil))
}
