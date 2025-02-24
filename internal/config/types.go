package config

import "time"

type AppConfig struct {
	HTTP HTTPConfig `envPrefix:"HTTP_"`
	DB   DBConfig   `envPrefix:"DB_"`

	BuildInfo BuildInfo
}

type HTTPConfig struct {
	Host            string        `env:"HOST" envDefault:"0.0.0.0"`
	Port            int           `env:"PORT" envDefault:"8080"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" envDefault:"120s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"20s"`
	CORS            struct {
		AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
		AllowedMethods []string `env:"ALLOWED_METHODS"`
		AllowedHeaders []string `env:"ALLOWED_HEADERS"`
	} `envPrefix:"CORS_"`
}

type DBConfig struct {
	Driver      string `env:"DRIVER" envDefault:"postgres"`
	Host        string `env:"HOST" envDefault:"127.0.0.1"`
	Port        int    `env:"PORT" envDefault:"5432"`
	User        string `env:"USER" envDefault:"postgres"`
	Password    string `env:"PASSWORD" envDefault:"postgres"`
	Name        string `env:"NAME" envDefault:"sandbox"`
	Schema      string `env:"SCHEMA" envDefault:"ebook"`
	MaxIdle     int    `env:"MAX_IDLE" envDefault:"2"`
	MaxOpen     int    `env:"MAX_OPEN" envDefault:"10"`
	SSLMode     string `env:"SSL_MODE" envDefault:"disable"`
	Timezone    string `env:"TIMEZONE" envDefault:"UTC"`
	AutoMigrate bool   `env:"AUTO_MIGRATE" envDefault:"false"`
}

type BuildInfo struct {
	Revision string
	Time     string
	Dirty    string
}
