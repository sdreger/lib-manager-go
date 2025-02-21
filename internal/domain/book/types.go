package book

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type Book struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Subtitle      string    `json:"subtitle"`
	Description   string    `json:"description"`
	ISBN10        string    `json:"isbn10"`
	ISBN13        int64     `json:"isbn13"`
	ASIN          string    `json:"asin"`
	Pages         uint16    `json:"pages"`
	PublisherURL  string    `json:"publisher_url"`
	Edition       uint8     `json:"edition"`
	PubDate       time.Time `json:"pub_date"`
	BookFileName  string    `json:"book_file_name"`
	BookFileSize  int64     `json:"book_file_size"`
	CoverFileName string    `json:"cover_file_name"`
	Language      string    `json:"language"`
	Publisher     string    `json:"publisher"`
	Authors       []string  `json:"authors"`
	Categories    []string  `json:"categories"`
	FileTypes     []string  `json:"file_types"`
	Tags          []string  `json:"tags"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type bookEntity struct {
	ID            int64          `db:"id"`
	Title         string         `db:"title"`
	Subtitle      sql.NullString `db:"subtitle"`
	Description   sql.NullString `db:"description"`
	ISBN10        sql.NullString `db:"isbn10"`
	ISBN13        sql.NullInt64  `db:"isbn13"`
	ASIN          sql.NullString `db:"asin"`
	Pages         uint16         `db:"pages"`
	PublisherURL  string         `db:"publisher_url"`
	Edition       uint8          `db:"edition"`
	PubDate       time.Time      `db:"pub_date"`
	BookFileName  string         `db:"book_file_name"`
	BookFileSize  int64          `db:"book_file_size"`
	CoverFileName string         `db:"cover_file_name"`
	Language      string         `db:"language"`
	Publisher     string         `db:"publisher"`
	Authors       pq.StringArray `db:"authors"`
	Categories    pq.StringArray `db:"categories"`
	FileTypes     pq.StringArray `db:"file_types"`
	Tags          pq.StringArray `db:"tags"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}

type LookupItem struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Subtitle      string    `json:"subtitle"`
	ISBN10        string    `json:"isbn10"`
	ISBN13        int64     `json:"isbn13"`
	ASIN          string    `json:"asin"`
	Pages         uint16    `json:"pages"`
	Edition       uint8     `json:"edition"`
	PubDate       time.Time `json:"pub_date"`
	BookFileSize  int64     `json:"book_file_size"`
	CoverFileName string    `json:"cover_file_name"`
	Publisher     string    `json:"publisher"`
	Language      string    `json:"language"`
	AuthorIDs     []int64   `json:"author_ids"`
	CategoryIDs   []int64   `json:"category_ids"`
	FileTypeIDs   []int64   `json:"file_types_ids"`
	TagIDs        []int64   `json:"tag_ids"`
}

type lookupEntity struct {
	ID            int64          `db:"id"`
	Title         string         `db:"title"`
	Subtitle      sql.NullString `db:"subtitle"`
	ISBN10        sql.NullString `db:"isbn10"`
	ISBN13        sql.NullInt64  `db:"isbn13"`
	ASIN          sql.NullString `db:"asin"`
	Pages         uint16         `db:"pages"`
	Edition       uint8          `db:"edition"`
	PubDate       time.Time      `db:"pub_date"`
	BookFileSize  int64          `db:"book_file_size"`
	CoverFileName string         `db:"cover_file_name"`
	Publisher     string         `db:"publisher"`
	Language      string         `db:"language"`
	AuthorIDs     pq.Int64Array  `db:"author_ids"`
	CategoryIDs   pq.Int64Array  `db:"category_ids"`
	FileTypeIDs   pq.Int64Array  `db:"file_types_ids"`
	TagIDs        pq.Int64Array  `db:"tag_ids"`
	Total         int64          `db:"total"`
}
