package book

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type DBStore struct {
	db *sqlx.DB
}

func NewDBStore(db *sqlx.DB) *DBStore {
	return &DBStore{db: db}
}

func (s *DBStore) GetByID(ctx context.Context, bookID int64) (Book, error) {
	var book bookEntity
	query := `SELECT books.id AS id, title, subtitle, description, isbn10, isbn13, asin,
       pages, publisher_url, edition, pub_date, book_file_name, book_file_size,
       cover_file_name, books.created_at AS created_at, books.updated_at AS updated_at,
       publishers.name                     AS publisher,
       languages.name                      AS language,
       ARRAY_AGG(DISTINCT authors.name)    AS authors,
       ARRAY_AGG(DISTINCT categories.name) AS categories,
       ARRAY_AGG(DISTINCT file_types.name) AS file_types,
       ARRAY_REMOVE(ARRAY_AGG(DISTINCT tags.name), NULL) AS tags
FROM ebook.books
         LEFT JOIN ebook.languages on books.language_id = languages.id
         LEFT JOIN ebook.publishers on books.publisher_id = publishers.id
         LEFT JOIN ebook.book_author on books.id = book_author.book_id
         LEFT JOIN ebook.authors on authors.id = book_author.author_id
         LEFT JOIN ebook.book_category ON books.id = book_category.book_id
         LEFT JOIN ebook.categories ON book_category.category_id = categories.id
         LEFT JOIN ebook.book_file_type ON books.id = book_file_type.book_id
         LEFT JOIN ebook.file_types ON book_file_type.file_type_id = file_types.id
         LEFT JOIN ebook.book_tag ON books.id = book_tag.book_id
         LEFT JOIN ebook.tags ON book_tag.tag_id = tags.id
WHERE books.id = $1
GROUP BY books.id, title, subtitle, description, isbn10, isbn13, asin, pages, publisher_url,
         edition, pub_date, book_file_name, book_file_size, cover_file_name, 
         books.created_at, books.updated_at, publishers.name, languages.name
`
	err := s.db.GetContext(ctx, &book, query, bookID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Book{}, ErrNotFound
		}

		return Book{}, err
	}

	return s.FromEntity(book), nil
}

func (s *DBStore) FromEntity(book bookEntity) Book {
	fullView := Book{
		ID:            book.ID,
		Title:         book.Title,
		Pages:         book.Pages,
		PublisherURL:  book.PublisherURL,
		Edition:       book.Edition,
		PubDate:       book.PubDate,
		BookFileName:  book.BookFileName,
		BookFileSize:  book.BookFileSize,
		CoverFileName: book.CoverFileName,
		Language:      book.Language,
		Publisher:     book.Publisher,
		Authors:       book.Authors,
		Categories:    book.Categories,
		FileTypes:     book.FileTypes,
		Tags:          book.Tags,
		CreatedAt:     book.CreatedAt,
		UpdatedAt:     book.UpdatedAt,
	}
	if book.Subtitle.Valid {
		fullView.Subtitle = book.Subtitle.String
	}
	if book.Description.Valid {
		fullView.Description = book.Description.String
	}
	if book.ISBN10.Valid {
		fullView.ISBN10 = book.ISBN10.String
	}
	if book.ISBN13.Valid {
		fullView.ISBN13 = book.ISBN13.Int64
	}
	if book.ASIN.Valid {
		fullView.ASIN = book.ASIN.String
	}
	return fullView
}
