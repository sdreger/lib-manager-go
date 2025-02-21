package book

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sdreger/lib-manager-go/internal/paging"
)

type DBStore struct {
	db *sqlx.DB
}

func NewDBStore(db *sqlx.DB) *DBStore {
	return &DBStore{db: db}
}

// GetByID - returns a book by its ID if present, otherwise returns ErrNotFound
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

	return s.fromEntity(book), nil
}

// Lookup - returns a paginated, sorted and filtered slice of lookup items
func (s *DBStore) Lookup(ctx context.Context, page paging.PageRequest, sort paging.Sort, filter Filter) (
	[]LookupItem, int64, error) {

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.Select(`books.id, title, subtitle, isbn10, isbn13, asin, pages, edition, pub_date, 
       			book_file_size, cover_file_name, publishers.name as publisher, languages.name as language,
			 	array_agg(DISTINCT ba.author_id) as author_ids,
			 	array_agg(DISTINCT bc.category_id) as category_ids,
       		    array_agg(DISTINCT bft.file_type_id) as file_types_ids,
		        array_remove(array_agg(DISTINCT bt.tag_id), NULL) as tag_ids,
				count(*) over() as total`).
		From("ebook.books").
		LeftJoin("ebook.publishers on books.publisher_id = publishers.id").
		LeftJoin("ebook.languages on books.language_id = languages.id").
		LeftJoin("ebook.book_author ba on books.id = ba.book_id").
		LeftJoin("ebook.book_file_type bft on books.id = bft.book_id").
		LeftJoin("ebook.book_category bc on books.id = bc.book_id").
		LeftJoin("ebook.book_tag bt on books.id = bt.book_id").
		GroupBy(`books.id, title, subtitle, isbn10, isbn13, asin, pages, 
			           pub_date, book_file_size, cover_file_name, publisher, language`).
		OrderBy(sort.GetOrderBy()).
		Limit(page.Limit()).
		Offset(page.Offset())

	// filter by ISBN10/ISBN13/ASIN overrides all other filters
	if filter.SBN != "" {
		query = query.Where(
			sq.Or{
				sq.Eq{"books.ISBN10": filter.SBN},
				sq.Eq{"books.ISBN13::varchar": filter.SBN},
				sq.Eq{"books.ASIN": filter.SBN},
			},
		)
	} else {
		if len(filter.Languages) > 0 {
			query = query.Where("language_id = ANY(?)", pq.Array(filter.Languages))
		}
		if len(filter.Publishers) > 0 {
			query = query.Where("publisher_id = ANY(?)", pq.Array(filter.Publishers))
		}
		if len(filter.Authors) > 0 {
			query = query.Having("(array_agg(ba.author_id) && ?)", pq.Array(filter.Authors))
		}
		if len(filter.Categories) > 0 {
			query = query.Having("(array_agg(bc.category_id) && ?)", pq.Array(filter.Categories))
		}
		if len(filter.FileTypes) > 0 {
			query = query.Having("(array_agg(bft.file_type_id) && ?)", pq.Array(filter.FileTypes))
		}

		if len(filter.Tags) > 0 {
			query = query.Having("(array_agg(bt.tag_id) && ?)", pq.Array(filter.Tags))
		}

		if len(filter.Query) > 0 {
			query = query.Where(sq.ILike{"books.title": "%" + filter.Query + "%"})
		}
	}

	sqlQuery, queryParams, err := query.ToSql()
	if err != nil {
		return nil, 0, err
	}

	var rows []lookupEntity
	err = s.db.SelectContext(ctx, &rows, sqlQuery, queryParams...)
	if err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if len(rows) > 0 {
		total = rows[0].Total
	}

	lookupItems := make([]LookupItem, len(rows))
	for i, row := range rows {
		lookupItems[i] = s.fromLookupEntity(row)
	}

	return lookupItems, total, nil
}

func (s *DBStore) fromEntity(book bookEntity) Book {
	result := Book{
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
		result.Subtitle = book.Subtitle.String
	}
	if book.Description.Valid {
		result.Description = book.Description.String
	}
	if book.ISBN10.Valid {
		result.ISBN10 = book.ISBN10.String
	}
	if book.ISBN13.Valid {
		result.ISBN13 = book.ISBN13.Int64
	}
	if book.ASIN.Valid {
		result.ASIN = book.ASIN.String
	}
	return result
}

func (s *DBStore) fromLookupEntity(book lookupEntity) LookupItem {
	result := LookupItem{
		ID:            book.ID,
		Title:         book.Title,
		Pages:         book.Pages,
		Edition:       book.Edition,
		PubDate:       book.PubDate,
		BookFileSize:  book.BookFileSize,
		CoverFileName: book.CoverFileName,
		Publisher:     book.Publisher,
		Language:      book.Language,
		AuthorIDs:     book.AuthorIDs,
		CategoryIDs:   book.CategoryIDs,
		FileTypeIDs:   book.FileTypeIDs,
		TagIDs:        book.TagIDs,
	}
	if book.Subtitle.Valid {
		result.Subtitle = book.Subtitle.String
	}
	if book.ISBN10.Valid {
		result.ISBN10 = book.ISBN10.String
	}
	if book.ISBN13.Valid {
		result.ISBN13 = book.ISBN13.Int64
	}
	if book.ASIN.Valid {
		result.ASIN = book.ASIN.String
	}

	return result
}
