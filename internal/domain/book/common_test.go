package book

import "time"

const (
	bookID            = int64(1)
	bookTitle         = "CockroachDB"
	bookSubtitle      = "The Definitive Guide"
	bookDescription   = "Get the lowdown on CockroachDB"
	bookISBN10        = "1234567890"
	bookISBN13        = 9781234567890
	bookASIN          = "BH34567890"
	bookPages         = 256
	bookEdition       = 2
	bookPublisherURL  = "https://amazon.com/dp/1234567890.html"
	bookPubDate       = "2022-07-19"
	bookFileName      = "OReilly.CockroachDB.2nd.Edition.1234567890.zip"
	bookFileSize      = 5192
	bookCoverFileName = "1234567890.jpg"
	bookLanguage      = "English"
	bookPublisher     = "OReilly"
	bookAuthorID01    = int64(1)
	bookAuthorID02    = int64(2)
	bookAuthor01      = "John Doe"
	bookAuthor02      = "Amanda Lee"
	bookCategoryID01  = int64(1)
	bookCategoryID02  = int64(2)
	bookCategoryID03  = int64(3)
	bookCategory01    = "Computer Science"
	bookCategory02    = "Computers"
	bookCategory03    = "Programming"
	bookFileTypeID01  = int64(1)
	bookFileTypeID02  = int64(2)
	bookFileType01    = "pdf"
	bookFileType02    = "epub"
	bookTagID01       = int64(1)
	bookTagID02       = int64(2)
	bookTag01         = "programming"
	bookTag02         = "database"
)

func getTestBook() Book {
	bookPubDate, _ := time.Parse(time.DateOnly, bookPubDate)
	createdAt := time.Date(2025, time.April, 10, 9, 15, 10, 0, time.UTC)
	updatedAt := time.Date(2025, time.April, 15, 10, 25, 15, 0, time.UTC)
	return Book{
		ID:            bookID,
		Title:         bookTitle,
		Subtitle:      bookSubtitle,
		Description:   bookDescription,
		ISBN10:        bookISBN10,
		ISBN13:        bookISBN13,
		ASIN:          bookASIN,
		Pages:         bookPages,
		PublisherURL:  bookPublisherURL,
		Edition:       bookEdition,
		PubDate:       bookPubDate,
		BookFileName:  bookFileName,
		BookFileSize:  bookFileSize,
		CoverFileName: bookCoverFileName,
		Language:      bookLanguage,
		Publisher:     bookPublisher,
		Authors:       []string{bookAuthor01, bookAuthor02},
		Categories:    []string{bookCategory01, bookCategory02, bookCategory03},
		FileTypes:     []string{bookFileType01, bookFileType02},
		Tags:          []string{bookTag01, bookTag02},
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func getTestLookupItem() LookupItem {
	bookPubDate, _ := time.Parse(time.DateOnly, bookPubDate)
	return LookupItem{
		ID:            bookID,
		Title:         bookTitle,
		Subtitle:      bookSubtitle,
		ISBN10:        bookISBN10,
		ISBN13:        bookISBN13,
		ASIN:          bookASIN,
		Pages:         bookPages,
		Edition:       bookEdition,
		PubDate:       bookPubDate,
		BookFileSize:  bookFileSize,
		CoverFileName: bookCoverFileName,
		Language:      bookLanguage,
		Publisher:     bookPublisher,
		AuthorIDs:     []int64{bookAuthorID01, bookAuthorID02},
		CategoryIDs:   []int64{bookCategoryID01, bookCategoryID02, bookCategoryID03},
		FileTypeIDs:   []int64{bookFileTypeID01, bookFileTypeID02},
		TagIDs:        []int64{bookTagID01, bookTagID02},
	}
}
