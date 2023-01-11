package config

import (
	"github.com/amirhnajafiz/fhs/internal/cmd/server"
	"github.com/amirhnajafiz/fhs/internal/storage"
	"github.com/amirhnajafiz/fhs/internal/telemetry/config"
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
		Storage: storage.Config{
			AccessKeyID:     "",
			SecretAccessKey: "",
			Region:          "",
			Bucket:          "",
			Endpoint:        "",
		},
		Server: server.Config{
			Host: "0.0.0.0",
			Port: 8080,
		},
	}
}
