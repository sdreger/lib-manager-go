package publisher

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sdreger/lib-manager-go/internal/paging"
	"github.com/sdreger/lib-manager-go/internal/tests"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"os"
	"testing"
)

type TestStoreSuite struct {
	suite.Suite
	db            *sqlx.DB
	testContainer *postgres.PostgresContainer
	store         *DBStore
}

func (s *TestStoreSuite) SetupSuite() {
	testContainer := tests.StartDBTestContainer(s.T())
	dbConfig := tests.GetTestDBConfig(s.T(), testContainer)
	connection := tests.SetUpTestDB(s.Suite.Require(), dbConfig, testContainer)

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
	publishersFound := 4
	s.Equal(int64(publishersFound), total)
	s.Len(response, publishersFound)
	s.Require().Less(response[0].Name, response[1].Name)
	s.Require().Less(response[1].Name, response[2].Name)
	s.Require().Less(response[2].Name, response[3].Name)
}

func (s *TestStoreSuite) Test_Lookup_OrderByName_OnePage() {
	requestValues := map[string][]string{"page": {"2"}, "size": {"2"}, "sort": {"id,desc"}}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err, "failed to perform lookup request")
	publishersFound := 2
	expectedTotal := int64(4)
	s.Equal(expectedTotal, total)
	s.Len(response, publishersFound)
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

	err := prepareTestData(s.testContainer, "testdata/publisher_lookup.sql")
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

	return tests.ExecSQL(testContainer, string(file))
}
