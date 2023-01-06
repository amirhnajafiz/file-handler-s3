package config

import (
	"fmt"
	"log"

	"github.com/amirhnajafiz/hls/internal/cmd/server"
	"github.com/amirhnajafiz/hls/internal/storage"
	"github.com/amirhnajafiz/hls/internal/telemetry/config"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
)

type Config struct {
	Telemetry config.Config  `koanf:"telemetry"`
	Storage   storage.Config `koanf:"storage"`
	Server    server.Config  `koanf:"server"`
}

func Load() Config {
	var instance Config

	k := koanf.New(".")

	// load default
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		_ = fmt.Errorf("error loading deafult: %v\n", err)
	}

	// load configs file
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		_ = fmt.Errorf("error loading config.yaml file: %v\n", err)
	}

	// unmarshalling
	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %v\n", err)
	}

	return instance
}
