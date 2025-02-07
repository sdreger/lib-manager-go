package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	defaultHost               = "0.0.0.0"
	defaultPort               = "8080"
	defaultReadTimeout, _     = time.ParseDuration("5s")
	defaultWriteTimeout, _    = time.ParseDuration("10s")
	defaultIdleTimeout, _     = time.ParseDuration("120s")
	defaultShutdownTimeout, _ = time.ParseDuration("20s")
)

func TestNewConfigDefaults(t *testing.T) {

	config, err := New()
	if assert.NoError(t, err, "should parse default config") {
		assert.NotNil(t, config, "config should not be nil")
		assert.Equal(t, defaultHost, config.HTTP.Host)
		assert.Equal(t, defaultPort, config.HTTP.Port)
		assert.Equal(t, defaultReadTimeout, config.HTTP.ReadTimeout)
		assert.Equal(t, defaultWriteTimeout, config.HTTP.WriteTimeout)
		assert.Equal(t, defaultIdleTimeout, config.HTTP.IdleTimeout)
		assert.Equal(t, defaultShutdownTimeout, config.HTTP.ShutdownTimeout)
	}
}

func TestNewConfigCustomEnv(t *testing.T) {
	customHost := "127.0.0.1"
	customPort := "9090"
	customReadTimeout, _ := time.ParseDuration("3s")
	customWriteTimeout, _ := time.ParseDuration("5s")
	customIdleTimeout, _ := time.ParseDuration("20s")
	customShutdownTimeout, _ := time.ParseDuration("30s")
	_ = os.Setenv(getEnvKey("HTTP_HOST"), customHost)
	_ = os.Setenv(getEnvKey("HTTP_PORT"), customPort)
	_ = os.Setenv(getEnvKey("HTTP_READ_TIMEOUT"), customReadTimeout.String())
	_ = os.Setenv(getEnvKey("HTTP_WRITE_TIMEOUT"), customWriteTimeout.String())
	_ = os.Setenv(getEnvKey("HTTP_IDLE_TIMEOUT"), customIdleTimeout.String())
	_ = os.Setenv(getEnvKey("HTTP_SHUTDOWN_TIMEOUT"), customShutdownTimeout.String())

	config, err := New()
	if assert.NoError(t, err, "should parse custom config") {
		assert.NotNil(t, config, "config should not be nil")
		assert.Equal(t, customHost, config.HTTP.Host)
		assert.Equal(t, customPort, config.HTTP.Port)
		assert.Equal(t, customReadTimeout, config.HTTP.ReadTimeout)
		assert.Equal(t, customWriteTimeout, config.HTTP.WriteTimeout)
		assert.Equal(t, customIdleTimeout, config.HTTP.IdleTimeout)
		assert.Equal(t, customShutdownTimeout, config.HTTP.ShutdownTimeout)
	}
}

func TestNewConfigWithEmptyEnv(t *testing.T) {
	_ = os.Setenv(getEnvKey("HTTP_HOST"), "")
	_ = os.Setenv(getEnvKey("HTTP_PORT"), "")
	_ = os.Setenv(getEnvKey("HTTP_READ_TIMEOUT"), "")
	_ = os.Setenv(getEnvKey("HTTP_WRITE_TIMEOUT"), "")
	_ = os.Setenv(getEnvKey("HTTP_IDLE_TIMEOUT"), "")
	_ = os.Setenv(getEnvKey("HTTP_SHUTDOWN_TIMEOUT"), "")

	config, err := New()
	if assert.NoError(t, err, "should parse default config") {
		assert.NotNil(t, config, "config should not be nil")
		assert.Equal(t, defaultHost, config.HTTP.Host)
		assert.Equal(t, defaultPort, config.HTTP.Port)
		assert.Equal(t, defaultReadTimeout, config.HTTP.ReadTimeout)
		assert.Equal(t, defaultWriteTimeout, config.HTTP.WriteTimeout)
		assert.Equal(t, defaultIdleTimeout, config.HTTP.IdleTimeout)
		assert.Equal(t, defaultShutdownTimeout, config.HTTP.ShutdownTimeout)
	}
}

func getEnvKey(key string) string {
	return envPrefix + key
}
