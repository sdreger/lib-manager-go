package book

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
	"strings"
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

func (s *TestStoreSuite) Test_GetByID_ErrorNotFound() {
	_, err := s.store.GetByID(context.Background(), bookID)
	s.ErrorIs(err, ErrNotFound)
}

func (s *TestStoreSuite) Test_GetByID_ErrorRelations() {
	ctx := context.Background()
	err := prepareTestData(s.testContainer, "testdata/book_required_relations.sql")
	s.Require().NoError(err, "failed to load test SQL file")

	// The book has no 'authors', which is error
	// Scan error on column index 18, name "authors": pq: parsing array element index 0: cannot convert nil to string
	_, err = s.store.GetByID(ctx, bookID)
	s.Require().Error(err)
}

func (s *TestStoreSuite) Test_GetByID_AllRelations() {
	ctx := context.Background()
	err := prepareTestData(s.testContainer, "testdata/book_all_relations.sql")
	s.Require().NoError(err, "failed to load test SQL file")

	response, err := s.store.GetByID(ctx, bookID)
	s.Require().NoError(err)
	testBook := getTestBook()
	s.Equal(testBook.ID, response.ID)
	s.Equal(testBook.Title, response.Title)
	s.Equal(testBook.Subtitle, response.Subtitle)
	s.Equal(testBook.Description, response.Description)
	s.Equal(testBook.ISBN10, response.ISBN10)
	s.Equal(testBook.ISBN13, response.ISBN13)
	s.Equal(testBook.ASIN, response.ASIN)
	s.Equal(testBook.Pages, response.Pages)
	s.Equal(testBook.PublisherURL, response.PublisherURL)
	s.Equal(testBook.Edition, response.Edition)
	s.Equal(testBook.PubDate, response.PubDate.In(time.UTC))
	s.Equal(testBook.BookFileName, response.BookFileName)
	s.Equal(testBook.BookFileSize, response.BookFileSize)
	s.Equal(testBook.CoverFileName, response.CoverFileName)
	s.Equal(testBook.Language, response.Language)
	s.Equal(testBook.Publisher, response.Publisher)
	s.ElementsMatch(testBook.Authors, response.Authors)
	s.ElementsMatch(testBook.Categories, response.Categories)
	s.ElementsMatch(testBook.FileTypes, response.FileTypes)
	s.ElementsMatch(testBook.Tags, response.Tags)
	now := time.Now()
	s.LessOrEqual(response.CreatedAt, now)
	s.LessOrEqual(response.UpdatedAt, now)
}

func (s *TestStoreSuite) Test_Lookup_NoFilters() {
	requestValues := map[string][]string{"page": {"1"}, "size": {"10"}}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err)
	booksFound := 3
	s.Equal(int64(booksFound), total)
	s.Len(response, booksFound)
}

func (s *TestStoreSuite) Test_Lookup_AllFilters() {
	query := "book 03"
	publisherID := "2"
	languageID := "2"
	authorID := "2"
	categoryID := "2"
	fileTypeID := "2"
	tagID := "2"

	authorIDNum, _ := strconv.Atoi(authorID)
	categoryIDNum, _ := strconv.Atoi(categoryID)
	fileTypeIDNum, _ := strconv.Atoi(fileTypeID)
	tagIDNum, _ := strconv.Atoi(tagID)

	requestValues := map[string][]string{
		"page":      {"1"},
		"size":      {"10"},
		"publisher": {publisherID},
		"language":  {languageID},
		"author":    {authorID},
		"category":  {categoryID},
		"file_type": {fileTypeID},
		"tag":       {tagID},
		"query":     {query},
	}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err)
	booksFound := 1
	s.Equal(int64(booksFound), total)
	s.Len(response, booksFound)
	book03 := response[0]
	s.Equal(int64(3), book03.ID)
	s.Equal("Manning", book03.Publisher)
	s.Equal("German", book03.Language)
	s.Contains(book03.AuthorIDs, int64(authorIDNum))
	s.Contains(book03.CategoryIDs, int64(categoryIDNum))
	s.Contains(book03.FileTypeIDs, int64(fileTypeIDNum))
	s.Contains(book03.TagIDs, int64(tagIDNum))
	s.Contains(strings.ToLower(book03.Title), strings.ToLower(query))
}

func (s *TestStoreSuite) Test_Lookup_PublisherFilters() {
	requestValues := map[string][]string{"page": {"1"}, "size": {"10"}, "publisher": {"1"}}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err)
	booksFound := 2
	s.Equal(int64(booksFound), total)
	s.Len(response, booksFound)
	book01 := response[0]
	book02 := response[1]
	s.Equal(int64(1), book01.ID)
	s.Equal(int64(2), book02.ID)
}

func (s *TestStoreSuite) Test_Lookup_LanguageFilters() {
	requestValues := map[string][]string{"page": {"1"}, "size": {"10"}, "language": {"1"}}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err)
	booksFound := 2
	s.Equal(int64(booksFound), total)
	s.Len(response, booksFound)
	book01 := response[0]
	book02 := response[1]
	s.Equal(int64(1), book01.ID)
	s.Equal(int64(2), book02.ID)
}

func (s *TestStoreSuite) Test_Lookup_AuthorFilters() {
	requestValues := map[string][]string{"page": {"1"}, "size": {"10"}, "author": {"1"}}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err)
	booksFound := 2
	s.Equal(int64(booksFound), total)
	s.Len(response, booksFound)
	book01 := response[0]
	book02 := response[1]
	s.Equal(int64(1), book01.ID)
	s.Equal(int64(2), book02.ID)
}

func (s *TestStoreSuite) Test_Lookup_CategoryFilters() {
	requestValues := map[string][]string{"page": {"1"}, "size": {"10"}, "category": {"1"}}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err)
	booksFound := 2
	s.Equal(int64(booksFound), total)
	s.Len(response, booksFound)
	book01 := response[0]
	book02 := response[1]
	s.Equal(int64(1), book01.ID)
	s.Equal(int64(2), book02.ID)
}

func (s *TestStoreSuite) Test_Lookup_FileTypeFilters() {
	requestValues := map[string][]string{"page": {"1"}, "size": {"10"}, "file_type": {"1"}}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err)
	booksFound := 2
	s.Equal(int64(booksFound), total)
	s.Len(response, booksFound)
	book01 := response[0]
	book02 := response[1]
	s.Equal(int64(1), book01.ID)
	s.Equal(int64(2), book02.ID)
}

func (s *TestStoreSuite) Test_Lookup_TagFilters() {
	requestValues := map[string][]string{"page": {"1"}, "size": {"10"}, "tag": {"1"}}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err)
	booksFound := 2
	s.Equal(int64(booksFound), total)
	s.Len(response, booksFound)
	book01 := response[0]
	book02 := response[1]
	s.Equal(int64(1), book01.ID)
	s.Equal(int64(2), book02.ID)
}

func (s *TestStoreSuite) Test_Lookup_QueryFilters() {
	requestValues := map[string][]string{"page": {"1"}, "size": {"10"}, "query": {"book 01"}}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err)
	booksFound := 1
	s.Equal(int64(booksFound), total)
	s.Len(response, booksFound)
	book01 := response[0]
	s.Equal(int64(1), book01.ID)
}

func (s *TestStoreSuite) Test_Lookup_Filters() {
	requestValues := map[string][]string{"page": {"1"}, "size": {"10"}, "sbn": {"3333333333"}}
	response, total, err := performLookupRequest(s, requestValues)
	s.Require().NoError(err)
	booksFound := 1
	s.Equal(int64(booksFound), total)
	s.Len(response, booksFound)
	book03 := response[0]
	s.Equal(int64(3), book03.ID)
}

func (s *TestStoreSuite) Test_Lookup_Error() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // should cause DB query error

	_, _, err := s.store.Lookup(ctx, paging.PageRequest{}, paging.Sort{}, Filter{})
	s.Require().Error(err, "lookup should fail")
}

func performLookupRequest(s *TestStoreSuite, requestValues map[string][]string) (
	[]LookupItem, int64, error) {

	err := prepareTestData(s.testContainer, "testdata/book_lookup_filter.sql")
	s.Require().NoError(err, "failed to load test SQL file")

	ctx := context.Background()
	pageRequest, err := paging.NewPageRequest(requestValues)
	s.Require().NoError(err, "failed to build page request")
	sort, err := paging.NewSort(requestValues, AllowedSortFields)
	s.Require().NoError(err, "failed to build sort")
	filter, err := NewFilter(requestValues)
	s.Require().NoError(err, "failed to build filter")

	return s.store.Lookup(ctx, pageRequest, sort, filter)
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
			wait.ForListeningPort("5432/tcp"),
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
