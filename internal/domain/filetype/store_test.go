package filetype

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/sdreger/lib-manager-go/internal/database"
	"github.com/sdreger/lib-manager-go/internal/paging"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
	dbUser     = "test"
	dBPassword = "test"
	dbName     = "test"
)

type TestStoreSuite struct {
	suite.Suite
	db            *sqlx.DB
	testContainer *postgres.PostgresContainer
	store         *DBStore
}

func (s *TestStoreSuite) SetupSuite() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	testContainer := startDBContainer(s.T())
	dbConfig := getTestDBConfig(s.T(), testContainer)
	connection, err := database.Open(dbConfig)
	s.Require().NoError(err, "failed to open database connection")
	err = database.Migrate(logger, dbConfig, connection.DB)
	s.Require().NoError(err, "failed to perform database migration")
	err = connection.Close() // required, to be able to create a testContainer snapshot
	s.Require().NoError(err, "failed to close database connection")

	err = testContainer.Snapshot(ctx)
	s.Require().NoError(err, "failed to create database snapshot")

	connection, err = database.Open(dbConfig)
	s.Require().NoError(err, "failed to open database connection")

	s.store = NewDBStore(connection)
	s.db = connection
	s.testContainer = testContainer
}

func (s *TestStoreSuite) SetupTest() {
	ctx := context.Background()
	err := s.testContainer.Restore(ctx)
	s.Require().NoError(err)
}

func (s *TestStoreSuite) TearDownSuite() {
	err := s.db.Close()
	s.Require().NoError(err, "failed to close database connection")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TestStoreSuite))
}

// -------------------- Tests --------------------

func (s *TestStoreSuite) Test_Lookup_OrderByName() {
	requestValues := map[string][]string{"page": {"1"}, "size": {"100"}, "sort": {"name,asc"}}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err, "failed to perform lookup request")
	booksFound := 4
	s.Equal(int64(booksFound), total)
	s.Len(response, booksFound)
	s.Require().Less(response[0].Name, response[1].Name)
	s.Require().Less(response[1].Name, response[2].Name)
	s.Require().Less(response[2].Name, response[3].Name)
}

func (s *TestStoreSuite) Test_Lookup_OrderByName_OnePage() {
	requestValues := map[string][]string{"page": {"2"}, "size": {"2"}, "sort": {"id,desc"}}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err, "failed to perform lookup request")
	booksFound := 2
	expectedTotal := int64(4)
	s.Equal(expectedTotal, total)
	s.Len(response, booksFound)
	s.Require().Greater(response[0].ID, response[1].ID)
}

func (s *TestStoreSuite) Test_Lookup_Error() {

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // should cause DB query error

	_, _, err := s.store.Lookup(ctx, paging.PageRequest{}, paging.Sort{})
	s.Require().Error(err, "lookup should fail")
}

func performLookupRequest(s *TestStoreSuite, requestValues map[string][]string) (
	[]LookupItem, int64, error) {

	err := prepareTestData(s.testContainer, "testdata/file_type_lookup.sql")
	s.Require().NoError(err, "failed to load test SQL file")

	ctx := context.Background()
	pageRequest, err := paging.NewPageRequest(requestValues)
	s.Require().NoError(err, "failed to build page request")
	sort, err := paging.NewSort(requestValues, AllowedSortFields)
	s.Require().NoError(err, "failed to build sort")

	return s.store.Lookup(ctx, pageRequest, sort)
}

func prepareTestData(testContainer *postgres.PostgresContainer, fileName string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	return execSQL(testContainer, string(file))
}

func execSQL(testContainer *postgres.PostgresContainer, statement string) error {
	ctx := context.Background()
	_, _, err := testContainer.Exec(ctx, []string{"psql", "-U", dbUser, "-d", dbName, "-c", statement})
	return err
}

func startDBContainer(t *testing.T) *postgres.PostgresContainer {
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
		),
	)
	testcontainers.CleanupContainer(t, pg)
	require.NoError(t, err, "error starting postgres test container")

	return pg
}

func getTestDBConfig(t *testing.T, pg *postgres.PostgresContainer) config.DBConfig {
	ctx := context.Background()
	appConfig, err := config.New()
	require.NoError(t, err, "failed to load application config")

	containerPort, err := pg.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err, "error getting mapped port 5432/tcp")
	port, err := strconv.Atoi(containerPort.Port())
	require.NoError(t, err, "error converting mapped port to int")

	dbConfig := appConfig.DB
	dbConfig.Port = port
	dbConfig.User = dbUser
	dbConfig.Password = dBPassword
	dbConfig.Name = dbName
	dbConfig.AutoMigrate = true

	return dbConfig
}
