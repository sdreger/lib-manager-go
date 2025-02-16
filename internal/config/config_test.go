package config

import (
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime/debug"
	"strconv"
	"testing"
	"time"
)

var (
	defaultHTTPHost               = "0.0.0.0"
	defaultHTTPPort               = 8080
	defaultHTTPReadTimeout, _     = time.ParseDuration("5s")
	defaultHTTPWriteTimeout, _    = time.ParseDuration("10s")
	defaultHTTPIdleTimeout, _     = time.ParseDuration("120s")
	defaultHTTPShutdownTimeout, _ = time.ParseDuration("20s")

	defaultDBDriver    = "postgres"
	defaultDBHost      = "127.0.0.1"
	defaultDBPort      = 5432
	defaultDBUser      = "postgres"
	defaultDBPassword  = "postgres"
	defaultDBName      = "sandbox"
	defaultDBSchema    = "ebook"
	defaultDBMaxIdle   = 2
	defaultDBMaxOpen   = 10
	defaultDBSSLMode   = "disable"
	defaultDBTimezone  = "UTC"
	defaultAutoMigrate = false
)

func TestNewConfigDefaults(t *testing.T) {

	config, err := New()
	if assert.NoError(t, err, "should parse default config") {
		if assert.NotEmpty(t, config.HTTP, "HTTP config should not be empty") {
			assert.Equal(t, defaultHTTPHost, config.HTTP.Host)
			assert.Equal(t, defaultHTTPPort, config.HTTP.Port)
			assert.Equal(t, defaultHTTPReadTimeout, config.HTTP.ReadTimeout)
			assert.Equal(t, defaultHTTPWriteTimeout, config.HTTP.WriteTimeout)
			assert.Equal(t, defaultHTTPIdleTimeout, config.HTTP.IdleTimeout)
			assert.Equal(t, defaultHTTPShutdownTimeout, config.HTTP.ShutdownTimeout)
		}

		if assert.NotEmpty(t, config.DB, "DB config should not be empty") {
			assert.Equal(t, defaultDBDriver, config.DB.Driver)
			assert.Equal(t, defaultDBHost, config.DB.Host)
			assert.Equal(t, defaultDBPort, config.DB.Port)
			assert.Equal(t, defaultDBUser, config.DB.User)
			assert.Equal(t, defaultDBPassword, config.DB.Password)
			assert.Equal(t, defaultDBName, config.DB.Name)
			assert.Equal(t, defaultDBSchema, config.DB.Schema)
			assert.Equal(t, defaultDBMaxIdle, config.DB.MaxIdle)
			assert.Equal(t, defaultDBMaxOpen, config.DB.MaxOpen)
			assert.Equal(t, defaultDBSSLMode, config.DB.SSLMode)
			assert.Equal(t, defaultDBTimezone, config.DB.Timezone)
			assert.Equal(t, defaultAutoMigrate, config.DB.AutoMigrate)
		}
	}
}

func TestNewConfigCustomHTTPEnv(t *testing.T) {
	customHost := "127.0.0.1"
	customPort := 9090
	customReadTimeout, _ := time.ParseDuration("3s")
	customWriteTimeout, _ := time.ParseDuration("5s")
	customIdleTimeout, _ := time.ParseDuration("20s")
	customShutdownTimeout, _ := time.ParseDuration("30s")
	_ = os.Setenv(getEnvKey("HTTP_HOST"), customHost)
	_ = os.Setenv(getEnvKey("HTTP_PORT"), strconv.Itoa(customPort))
	_ = os.Setenv(getEnvKey("HTTP_READ_TIMEOUT"), customReadTimeout.String())
	_ = os.Setenv(getEnvKey("HTTP_WRITE_TIMEOUT"), customWriteTimeout.String())
	_ = os.Setenv(getEnvKey("HTTP_IDLE_TIMEOUT"), customIdleTimeout.String())
	_ = os.Setenv(getEnvKey("HTTP_SHUTDOWN_TIMEOUT"), customShutdownTimeout.String())

	config, err := New()
	if assert.NoError(t, err, "should parse custom config") {
		assert.NotEmpty(t, config, "config should not be empty")
		assert.Equal(t, customHost, config.HTTP.Host)
		assert.Equal(t, customPort, config.HTTP.Port)
		assert.Equal(t, customReadTimeout, config.HTTP.ReadTimeout)
		assert.Equal(t, customWriteTimeout, config.HTTP.WriteTimeout)
		assert.Equal(t, customIdleTimeout, config.HTTP.IdleTimeout)
		assert.Equal(t, customShutdownTimeout, config.HTTP.ShutdownTimeout)
	}
}

func TestNewConfigCustomDBEnv(t *testing.T) {
	customDriver := "sqlite3"
	customHost := "192.168.1.100"
	customPort := 5555
	customUser := "testUser"
	customPassword := "testPassword"
	customName := "testDB"
	customSchema := "testSchema"
	customMaxIdle := 20
	customMaxOpen := 50
	customSSLMode := "enable"
	customTimezone := "CEST"
	customAutoMigrate := true

	_ = os.Setenv(getEnvKey("DB_DRIVER"), customDriver)
	_ = os.Setenv(getEnvKey("DB_HOST"), customHost)
	_ = os.Setenv(getEnvKey("DB_PORT"), strconv.Itoa(customPort))
	_ = os.Setenv(getEnvKey("DB_USER"), customUser)
	_ = os.Setenv(getEnvKey("DB_PASSWORD"), customPassword)
	_ = os.Setenv(getEnvKey("DB_NAME"), customName)
	_ = os.Setenv(getEnvKey("DB_SCHEMA"), customSchema)
	_ = os.Setenv(getEnvKey("DB_MAX_IDLE"), strconv.Itoa(customMaxIdle))
	_ = os.Setenv(getEnvKey("DB_MAX_OPEN"), strconv.Itoa(customMaxOpen))
	_ = os.Setenv(getEnvKey("DB_SSL_MODE"), customSSLMode)
	_ = os.Setenv(getEnvKey("DB_TIMEZONE"), customTimezone)
	_ = os.Setenv(getEnvKey("DB_AUTO_MIGRATE"), strconv.FormatBool(customAutoMigrate))

	config, err := New()
	if assert.NoError(t, err, "should parse custom config") {
		assert.NotEmpty(t, config, "config should not be empty")
		assert.Equal(t, customDriver, config.DB.Driver)
		assert.Equal(t, customHost, config.DB.Host)
		assert.Equal(t, customPort, config.DB.Port)
		assert.Equal(t, customUser, config.DB.User)
		assert.Equal(t, customPassword, config.DB.Password)
		assert.Equal(t, customName, config.DB.Name)
		assert.Equal(t, customSchema, config.DB.Schema)
		assert.Equal(t, customMaxIdle, config.DB.MaxIdle)
		assert.Equal(t, customMaxOpen, config.DB.MaxOpen)
		assert.Equal(t, customSSLMode, config.DB.SSLMode)
		assert.Equal(t, customTimezone, config.DB.Timezone)
		assert.Equal(t, customAutoMigrate, config.DB.AutoMigrate)
	}
}

func TestNewConfigWithEmptyEnv(t *testing.T) {
	_ = os.Setenv(getEnvKey("HTTP_HOST"), "")
	_ = os.Setenv(getEnvKey("HTTP_PORT"), "")
	_ = os.Setenv(getEnvKey("HTTP_READ_TIMEOUT"), "")
	_ = os.Setenv(getEnvKey("HTTP_WRITE_TIMEOUT"), "")
	_ = os.Setenv(getEnvKey("HTTP_IDLE_TIMEOUT"), "")
	_ = os.Setenv(getEnvKey("HTTP_SHUTDOWN_TIMEOUT"), "")

	_ = os.Setenv(getEnvKey("DB_Driver"), "")
	_ = os.Setenv(getEnvKey("DB_HOST"), "")
	_ = os.Setenv(getEnvKey("DB_PORT"), "")
	_ = os.Setenv(getEnvKey("DB_USER"), "")
	_ = os.Setenv(getEnvKey("DB_PASSWORD"), "")
	_ = os.Setenv(getEnvKey("DB_NAME"), "")
	_ = os.Setenv(getEnvKey("DB_SCHEMA"), "")
	_ = os.Setenv(getEnvKey("DB_MAX_IDLE"), "")
	_ = os.Setenv(getEnvKey("DB_MAX_OPEN"), "")
	_ = os.Setenv(getEnvKey("DB_SSL_MODE"), "")
	_ = os.Setenv(getEnvKey("DB_TIMEZONE"), "")
	_ = os.Setenv(getEnvKey("DB_AUTO_MIGRATE"), "")

	config, err := New()
	if assert.NoError(t, err, "should parse default config") {
		if assert.NotEmpty(t, config.HTTP, "HTTP config should not be empty") {
			http := config.HTTP
			assert.Equal(t, defaultHTTPHost, http.Host)
			assert.Equal(t, defaultHTTPPort, http.Port)
			assert.Equal(t, defaultHTTPReadTimeout, http.ReadTimeout)
			assert.Equal(t, defaultHTTPWriteTimeout, http.WriteTimeout)
			assert.Equal(t, defaultHTTPIdleTimeout, http.IdleTimeout)
			assert.Equal(t, defaultHTTPShutdownTimeout, http.ShutdownTimeout)
		}

		if assert.NotEmpty(t, config.DB, "DB config should not be empty") {
			assert.Equal(t, defaultDBHost, config.DB.Host)
			assert.Equal(t, defaultDBPort, config.DB.Port)
			assert.Equal(t, defaultDBUser, config.DB.User)
			assert.Equal(t, defaultDBPassword, config.DB.Password)
			assert.Equal(t, defaultDBName, config.DB.Name)
			assert.Equal(t, defaultDBSchema, config.DB.Schema)
			assert.Equal(t, defaultDBMaxIdle, config.DB.MaxIdle)
			assert.Equal(t, defaultDBMaxOpen, config.DB.MaxOpen)
			assert.Equal(t, defaultDBSSLMode, config.DB.SSLMode)
			assert.Equal(t, defaultDBTimezone, config.DB.Timezone)
			assert.Equal(t, defaultAutoMigrate, config.DB.AutoMigrate)
		}
	}
}

func TestNewConfigWithBuildInfo(t *testing.T) {
	vcsRevision := "2e509ddfc"
	vcsTime := "2025-02-07T07:22:13Z"
	vcsModified := "true"

	oldBuildInfoFunc := buildInfoFunc
	defer func() { buildInfoFunc = oldBuildInfoFunc }()
	buildInfoFunc = func() (info *debug.BuildInfo, ok bool) {
		return &debug.BuildInfo{
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: vcsRevision},
				{Key: "vcs.time", Value: vcsTime},
				{Key: "vcs.modified", Value: vcsModified},
			},
		}, true
	}

	config, err := New()
	if assert.NoError(t, err, "should parse default config") {
		assert.NotEmpty(t, config, "config should not be empty")
		buildInfo := config.BuildInfo
		if assert.NotNil(t, buildInfo, "build info should not be nil") {
			assert.Equal(t, buildInfo.Revision, vcsRevision, "revision should match")
			assert.Equal(t, buildInfo.Time, vcsTime, "time should match")
			assert.Equal(t, buildInfo.Dirty, vcsModified, "modified should match")
		}
	}
}

func TestNewConfigWithWrongPort(t *testing.T) {
	_ = os.Setenv(getEnvKey("HTTP_PORT"), "wrong value")
	config, err := New()
	if assert.ErrorIs(t, err, env.ParseError{}) {
		var parseError env.ParseError
		errors.As(err, &parseError)
		assert.Contains(t, parseError.Error(), `parsing "wrong value": invalid syntax`)
		assert.Empty(t, config, "config should not be empty")
	}
}

func getEnvKey(key string) string {
	return envPrefix + key
}
