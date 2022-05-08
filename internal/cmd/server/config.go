package server

type Config struct {
	Host string `koanf:"address"`
	Port int    `koanf:"port"`
}
