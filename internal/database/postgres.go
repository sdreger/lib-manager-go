package database

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sdreger/lib-manager-go/internal/config"
	"log/slog"
	"net/url"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const (
	scheme        = "postgres"
	timezoneKey   = "timezone"
	sslModeKey    = "sslmode"
	searchPathKey = "search_path"
)

// Open - opens a DB connection using provided configuration
func Open(config config.DBConfig) (*sqlx.DB, error) {
	urlValues := make(url.Values)
	urlValues.Set(timezoneKey, config.Timezone)
	urlValues.Set(sslModeKey, config.SSLMode)
	urlValues.Set(searchPathKey, config.Schema)

	dbURL := url.URL{
		Scheme:   scheme,
		Host:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		User:     url.UserPassword(config.User, config.Password),
		Path:     config.Name,
		RawQuery: urlValues.Encode(),
	}
	db, err := sqlx.Open(config.Driver, dbURL.String())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxOpen)
	db.SetMaxIdleConns(config.MaxIdle)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate - migrates the DB schema to the latest revision, if enabled in the DB config
func Migrate(logger *slog.Logger, dbConfig config.DBConfig, db *sql.DB) error {
	if dbConfig.AutoMigrate {
		logger.Info("starting database migration process")
		goose.SetBaseFS(embedMigrations)
		goose.SetLogger(slog.NewLogLogger(logger.Handler(), slog.LevelInfo))
		if err := goose.SetDialect(dbConfig.Driver); err != nil {
			return err
		}
		// if schema is created using the migration script, then 'goose_db_version' table
		// is created in the 'public' schema, not in the target one
		_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", dbConfig.Schema))
		if err != nil {
			return err
		}
		goose.SetTableName(dbConfig.Schema + ".goose_db_version")
		return goose.Up(db, "migrations")
	}

	logger.Info("automated database migration is disabled")
	return nil
}
