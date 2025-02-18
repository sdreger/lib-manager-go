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
