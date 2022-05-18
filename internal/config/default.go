package config

import (
	"github.com/amirhnajafiz/hls/internal/cmd/server"
	"github.com/amirhnajafiz/hls/internal/telemetry/config"
)

func Default() Config {
	return Config{
		Telemetry: config.Config{
			Trace: config.Trace{
				Enabled: false,
				Agent: config.Agent{
					Host: "",
					Port: "",
				},
			},
			Metric: config.Metric{
				Enabled: true,
				Address: ":1224",
			},
		},
		Server: server.Config{
			Host: "0.0.0.0",
			Port: 8080,
		},
	}
}
