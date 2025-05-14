//go:build !build

package tests

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/sdreger/lib-manager-go/internal/database"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log/slog"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	dbSnapshotTimeout = 5 * time.Second
	dbUser            = "test"
	dBPassword        = "test"
	dbName            = "test"
)

func StartDBTestContainer(t *testing.T) *postgres.PostgresContainer {
	pg, err := postgres.Run(context.Background(),
		"postgres:17.2-alpine3.21",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dBPassword),
		testcontainers.WithLogConsumers(&testcontainers.StdoutLogConsumer{}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
			wait.ForListeningPort("5432/tcp"),
		),
	)
	testcontainers.CleanupContainer(t, pg)
	require.NoError(t, err, "error starting postgres test container")

	return pg
}

func GetTestDBConfig(t *testing.T, pg *postgres.PostgresContainer) config.DBConfig {
	ctx := context.Background()
	appConfig, err := config.New()
	require.NoError(t, err, "failed to load application config")

	containerPort, err := pg.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err, "error getting mapped port 5432/tcp")
	port, err := strconv.Atoi(containerPort.Port())
	require.NoError(t, err, "error converting mapped port to int")
	host, err := pg.Host(ctx)
	require.NoError(t, err, "error getting test container host")

	dbConfig := appConfig.DB
	dbConfig.Host = host
	dbConfig.Port = port
	dbConfig.User = dbUser
	dbConfig.Password = dBPassword
	dbConfig.Name = dbName
	dbConfig.AutoMigrate = true

	return dbConfig
}

func SetUpTestDB(s *require.Assertions, dbConfig config.DBConfig, testContainer *postgres.PostgresContainer) *sqlx.DB {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))
	connection, err := database.Open(dbConfig)
	s.NoError(err, "failed to open database connection")

	err = database.Migrate(logger, dbConfig, connection.DB)
	s.NoError(err, "failed to perform database migration")

	err = connection.Close() // required, to be able to create a testContainer snapshot
	s.NoError(err, "failed to close database connection")

	snapshotCtx, snapshotCancelFunc := context.WithTimeout(context.Background(), dbSnapshotTimeout)
	defer snapshotCancelFunc()
	err = testContainer.Snapshot(snapshotCtx)
	s.NoError(err, "failed to create database snapshot")

	connection, err = database.Open(dbConfig)
	s.NoError(err, "failed to open database connection")

	return connection
}

func ExecSQL(testContainer *postgres.PostgresContainer, statement string) error {
	ctx := context.Background()
	_, _, err := testContainer.Exec(ctx, []string{"psql", "-U", dbUser, "-d", dbName, "-c", statement})
	return err
}
