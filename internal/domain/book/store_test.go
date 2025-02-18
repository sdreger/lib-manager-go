package book

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"path/filepath"
	"testing"
	"time"
)

const (
	dbUser     = "test"
	dBPassword = "test"
	dbName     = "test"
	schemaName = "ebook"
)

func TestStore(t *testing.T) {
	ctx := context.Background()

	testContainer := startDBContainer(t)
	execSQL(t, testContainer, "CREATE SCHEMA IF NOT EXISTS "+schemaName)
	runDBMigrations(t, testContainer)

	err := testContainer.Snapshot(ctx)
	require.NoError(t, err)

	connection := getDBConnection(t, testContainer)
	defer connection.Close()
	store := NewDBStore(connection)

	t.Run("getByIDErrorNotFound", func(t *testing.T) {
		t.Cleanup(func() {
			err := testContainer.Restore(ctx)
			require.NoError(t, err)
		})
		getByIDErrorNotFound(t, store)
	})

	t.Run("getByIDErrorRelations", func(t *testing.T) {
		t.Cleanup(func() {
			err := testContainer.Restore(ctx)
			require.NoError(t, err)
		})
		getByIDErrorRelations(t, store, testContainer)
	})

	t.Run("GetByIDBookWithAllRelations", func(t *testing.T) {
		t.Cleanup(func() {
			err := testContainer.Restore(ctx)
			require.NoError(t, err)
		})
		getByIDBookWithAllRelations(t, store, testContainer)
	})
}

func getByIDErrorNotFound(t *testing.T, store *DBStore) {
	_, err := store.GetByID(context.Background(), 1)
	assert.ErrorIs(t, err, ErrNotFound)
}

func getByIDErrorRelations(t *testing.T, store *DBStore, testContainer *postgres.PostgresContainer) {
	ctx := context.Background()
	addTestBook(t, testContainer)

	// The book has no 'authors', which is error
	// Scan error on column index 18, name "authors": pq: parsing array element index 0: cannot convert nil to string
	_, err := store.GetByID(ctx, 1)
	require.Error(t, err)
}

func getByIDBookWithAllRelations(t *testing.T, store *DBStore, testContainer *postgres.PostgresContainer) {
	ctx := context.Background()
	addTestBook(t, testContainer)
	addTestBookRelations(t, testContainer)

	response, err := store.GetByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "CockroachDB", response.Title)
}

func addTestBook(t *testing.T, testContainer *postgres.PostgresContainer) {
	publisherStmt := `INSERT INTO ebook.publishers (id, name) VALUES (1, 'OReilly');`
	languageStmt := `INSERT INTO ebook.languages (id, name) VALUES (1, 'English');`
	bookStmt := `INSERT INTO ebook.books (id, title, subtitle, description, isbn10, isbn13, asin, pages, edition, 
                 language_id, publisher_id, publisher_url, pub_date, book_file_name, book_file_size, cover_file_name)
				 VALUES (1, 'CockroachDB', 'The Definitive Guide', 'Get the lowdown on CockroachDB', '1234567890',
				 9781234567890, 'BH12345678', 256, 2, 1, 1, 'https://amazon.com/1234567890.html', '2022-07-19', 
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

func getDBConnection(t *testing.T, pg *postgres.PostgresContainer) *sqlx.DB {
	connectionArgs := fmt.Sprintf("search_path=%s&sslmode=disable&timezone=UTC", schemaName)
	connectionString, _ := pg.ConnectionString(context.Background(), connectionArgs)
	dbConnection, err := sqlx.Open("postgres", connectionString)
	require.NoError(t, err, "error connecting to postgres test container")

	return dbConnection
}

func runDBMigrations(t *testing.T, testContainer *postgres.PostgresContainer) {
	connection := getDBConnection(t, testContainer)
	absPath, err := filepath.Abs("../../database/migrations")
	require.NoError(t, err)

	err = goose.SetDialect("postgres")
	require.NoError(t, err)

	goose.SetTableName("ebook.goose_db_version")
	err = goose.Up(connection.DB, absPath)
	require.NoError(t, err)
	connection.Close() // required, to be able to create a testContainer snapshot
}
