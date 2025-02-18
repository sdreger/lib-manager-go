package book

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/sdreger/lib-manager-go/internal/database"
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
	store         Store
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
	addTestBook(s.T(), s.testContainer)

	// The book has no 'authors', which is error
	// Scan error on column index 18, name "authors": pq: parsing array element index 0: cannot convert nil to string
	_, err := s.store.GetByID(ctx, bookID)
	s.Require().Error(err)
}

func (s *TestStoreSuite) Test_GetByID_AllRelations() {
	ctx := context.Background()
	addTestBook(s.T(), s.testContainer)
	addTestBookRelations(s.T(), s.testContainer)

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

func addTestBook(t *testing.T, testContainer *postgres.PostgresContainer) {
	publisherStmt := `INSERT INTO ebook.publishers (id, name) VALUES (1, 'OReilly');`
	languageStmt := `INSERT INTO ebook.languages (id, name) VALUES (1, 'English');`
	bookStmt := `INSERT INTO ebook.books (id, title, subtitle, description, isbn10, isbn13, asin, pages, edition, 
                 language_id, publisher_id, publisher_url, pub_date, book_file_name, book_file_size, cover_file_name)
				 VALUES (1, 'CockroachDB', 'The Definitive Guide', 'Get the lowdown on CockroachDB', '1234567890',
				 9781234567890, 'BH34567890', 256, 2, 1, 1, 'https://amazon.com/dp/1234567890.html', '2022-07-19', 
				 'OReilly.CockroachDB.2nd.Edition.1234567890.zip', 5192, '1234567890.jpg');
`
	for _, stmt := range []string{publisherStmt, languageStmt, bookStmt} {
		execSQL(t, testContainer, stmt)
	}
}

func addTestBookRelations(t *testing.T, testContainer *postgres.PostgresContainer) {
	authorsStmt := `INSERT INTO ebook.authors (id, name) VALUES (1, 'John Doe'), (2, 'Amanda Lee');`
	bookAuthorsStmt := `INSERT INTO ebook.book_author (book_id, author_id) VALUES (1, 1), (1, 2);`
	categoriesStmt := `INSERT INTO ebook.categories (id, name, parent_id)
					   VALUES (1, 'Computer Science', null), (2, 'Computers', 1), (3, 'Programming', 2);`
	bookCategoriesStmt := `INSERT INTO ebook.book_category (book_id, category_id) VALUES (1, 1), (1, 2), (1, 3);`
	fileTypesStmt := `INSERT INTO ebook.file_types (id, name) VALUES (1, 'pdf'), (2, 'epub');`
	bookFileTypesStmt := `INSERT INTO ebook.book_file_type (book_id, file_type_id) VALUES (1, 1), (1, 2);`
	tagsStmt := `INSERT INTO ebook.tags (id, name) VALUES (1, 'programming'), (2, 'database');`
	bookTagsStmt := `INSERT INTO ebook.book_tag (book_id, tag_id) VALUES (1, 1), (1, 2);`

	for _, stmt := range []string{
		authorsStmt, bookAuthorsStmt, categoriesStmt, bookCategoriesStmt,
		fileTypesStmt, bookFileTypesStmt, tagsStmt, bookTagsStmt,
	} {
		execSQL(t, testContainer, stmt)
	}
}

func execSQL(t *testing.T, testContainer *postgres.PostgresContainer, statement string) {
	ctx := context.Background()
	_, _, err := testContainer.Exec(ctx, []string{"psql", "-U", dbUser, "-d", dbName, "-c", statement})
	require.NoError(t, err)
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
