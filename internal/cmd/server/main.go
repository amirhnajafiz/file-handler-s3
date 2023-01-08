package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/amirhnajafiz/hls/internal/port/http/handler"
	"github.com/amirhnajafiz/hls/internal/telemetry/metric"
	"go.opentelemetry.io/otel/trace"
)

// configure the songs' directory name
const mainDir = "files"

func New(cfg Config, t trace.Tracer, m metric.Metrics) {
	if _, err := os.Stat(mainDir); err != nil {
		if os.IsNotExist(err) {
			// file does not exist
			_ = os.Mkdir(mainDir, os.ModePerm)

			log.Println("Dir created")
		} else {
			log.Println("Dir exists")
		}
	}

	// handler init
	h := handler.Handler{
		Trace:  t,
		Metric: m,
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
	log.Printf("Serving %s on HTTP port: %v\n", mainDir, cfg.Port)

	// serve and log errors
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), nil))
}
