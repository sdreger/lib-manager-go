package config

import "time"

type AppConfig struct {
	HTTP HTTPConfig `envPrefix:"HTTP_"`
}

type HTTPConfig struct {
	Host            string        `env:"HOST" envDefault:"0.0.0.0"`
	Port            string        `env:"PORT" envDefault:"8080"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" envDefault:"120s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"20s"`
}
