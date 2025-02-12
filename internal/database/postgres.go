package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sdreger/lib-manager-go/internal/config"
	"net/url"
)

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
