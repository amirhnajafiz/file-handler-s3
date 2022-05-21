package cmd

import (
	"github.com/amirhnajafiz/hls/internal/cmd/server"
	"github.com/amirhnajafiz/hls/internal/config"
	"github.com/amirhnajafiz/hls/internal/telemetry/metric"
	"github.com/amirhnajafiz/hls/internal/telemetry/trace"
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

	// start server
	server.New(cfg.Server, t, m)
}
