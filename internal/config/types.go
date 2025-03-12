package config

import "time"

type AppConfig struct {
	HTTP      HTTPConfig      `envPrefix:"HTTP_"`
	DB        DBConfig        `envPrefix:"DB_"`
	BLOBStore BLOBStoreConfig `envPrefix:"BLOB_STORE_"`

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
	Driver                  string `env:"DRIVER" envDefault:"postgres"`
	Host                    string `env:"HOST" envDefault:"127.0.0.1"`
	Port                    int    `env:"PORT" envDefault:"5432"`
	User                    string `env:"USER" envDefault:"postgres"`
	Password                string `env:"PASSWORD" envDefault:"postgres"`
	Name                    string `env:"NAME" envDefault:"sandbox"`
	Schema                  string `env:"SCHEMA" envDefault:"ebook"`
	MaxIdle                 int    `env:"MAX_IDLE" envDefault:"2"`
	MaxOpen                 int    `env:"MAX_OPEN" envDefault:"10"`
	SSLMode                 string `env:"SSL_MODE" envDefault:"disable"`
	Timezone                string `env:"TIMEZONE" envDefault:"UTC"`
	AutoMigrate             bool   `env:"AUTO_MIGRATE" envDefault:"false"`
	MigrationLockTimeoutSec uint64 `env:"MIGRATION_LOCK_TIMEOUT_SEC" envDefault:"300"`
}

type BLOBStoreConfig struct {
	BookCoverBucket      string `env:"BOOK_COVER_BUCKET" envDefault:"ebook-covers"`
	MinioEndpoint        string `env:"MINIO_ENDPOINT" envDefault:"127.0.0.1:9000"`
	MinioAccessKeyID     string `env:"MINIO_ACCESS_KEY_ID" envDefault:"minio-access-key"`
	MinioSecretAccessKey string `env:"MINIO_ACCESS_SECRET_KEY" envDefault:"minio-secret-key"`
	MinioUseSSL          bool   `env:"MINIO_USE_SSL" envDefault:"false"`
}

type BuildInfo struct {
	Revision string
	Time     string
	Dirty    string
}
