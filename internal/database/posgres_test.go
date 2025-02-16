package database

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"strconv"
	"testing"
	"time"
)

func TestOpenDBConnection(t *testing.T) {
	ctx := context.Background()
	appConfig, err := config.New()
	require.NoError(t, err, "error loading config")
	dbConfig := appConfig.DB

	// https://golang.testcontainers.org/modules/postgres/
	// Postgres Docker container: https://hub.docker.com/_/postgres
	pg, err := postgres.Run(ctx,
		"postgres:17.2-alpine3.21",
		postgres.WithDatabase(dbConfig.Name),
		postgres.WithUsername(dbConfig.User),
		postgres.WithPassword(dbConfig.Password),
		testcontainers.WithLogConsumers(&testcontainers.StdoutLogConsumer{}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	require.NoError(t, err, "error starting postgres test container")
	defer func() {
		if err := testcontainers.TerminateContainer(pg); err != nil {
			t.Errorf("error terminating postgres test container: %v", err)
		}
	}()

	dbConnection, err := openDBConnection(t, dbConfig, pg)
	require.NoError(t, err, "error connecting to postgres test container")
	require.NotNil(t, dbConnection, "test container connection is nil")
	dbConnection.Close()
}

func openDBConnection(t *testing.T, dbConfig config.DBConfig, pg *postgres.PostgresContainer) (*sqlx.DB, error) {
	ctx := context.Background()
	containerPort, err := pg.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err, "error getting mapped port 5432/tcp")
	port, err := strconv.Atoi(containerPort.Port())
	require.NoError(t, err, "error converting mapped port to int")

	host, err := pg.Host(ctx)
	require.NoError(t, err, "error getting test container host")

	dbConfig.Port = port
	dbConfig.Host = host
	return Open(dbConfig)
}

func TestCannotPingDBConnection(t *testing.T) {
	appConfig, err := config.New()
	require.NoError(t, err, "error loading config")
	dbConfig := appConfig.DB

	conn, err := Open(dbConfig)
	require.Error(t, err, "should not be able to open postgres connection")
	assert.Nil(t, conn, "connection should be nil")
}

func TestCannotOpenDBConnection(t *testing.T) {
	appConfig, err := config.New()
	require.NoError(t, err, "error loading config")
	dbConfig := appConfig.DB
	dbConfig.Driver = "unknown"

	conn, err := Open(dbConfig)
	require.Error(t, err, "should not be able to open postgres connection")
	assert.Nil(t, conn, "connection should be nil")
}
