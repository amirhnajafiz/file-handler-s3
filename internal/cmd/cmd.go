package cmd

import (
	"github.com/amirhnajafiz/fhs/internal/cmd/server"
	"github.com/amirhnajafiz/fhs/internal/config"
	"github.com/amirhnajafiz/fhs/internal/port/s3"
	"github.com/amirhnajafiz/fhs/internal/telemetry/metric"
	"github.com/amirhnajafiz/fhs/internal/telemetry/trace"
)

func Execute() {
	// config load
	cfg := config.Load()

	// tracer init
	t := trace.New(cfg.Telemetry.Trace)

	// metric server init
	metric.NewServer(cfg.Telemetry.Metric).Start()

	// load metrics
	m := metric.NewMetrics()

	// creating storage handler
	s3handler, err := s3.New(cfg.Storage)
	if err != nil {
		panic(err)
	}

	// start server
	server.New(cfg.Server, t, m, s3handler)
}
