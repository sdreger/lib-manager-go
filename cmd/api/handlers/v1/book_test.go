package v1

import (
	"context"
	"encoding/json"
	"errors"
	apiErrors "github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/domain/book"
	"github.com/sdreger/lib-manager-go/internal/paging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

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

func TestBookHandler_RegisterBookHandler(t *testing.T) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))
	testRegistrar := handlers.RouteRegistrarMock{}
	h := NewBookHandler(logger, nil)
	h.RegisterHandler(&testRegistrar)

	assert.True(t, testRegistrar.IsRouteRegistered("GET /v1/books", h.GetBooks))
	assert.True(t, testRegistrar.IsRouteRegistered("GET /v1/books/{bookID}", h.GetBook))
}

func TestBookHandler_GetBook_Success(t *testing.T) {
	ctx := context.Background()
	handler := getBookHandler()

	testBook := getTestBook()

	mockService := NewMockBookService(t)
	testBookID := int64(1)
	mockService.EXPECT().GetBookByID(ctx, testBookID).Return(testBook, nil)
	injectBookMocks(handler, mockService)

	request := httptest.NewRequest("GET", "/v1/books/1", nil)
	request.SetPathValue("bookID", strconv.Itoa(int(testBookID)))
	recorder := httptest.NewRecorder()
	err := handler.GetBook(ctx, recorder, request)
	require.NoError(t, err, "should get a book")

	result := recorder.Result()
	defer result.Body.Close()
	require.Equal(t, http.StatusOK, result.StatusCode, "should get a 200 OK response")

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err, "should read body")
	var bookJSON map[string]book.Book
	_ = json.Unmarshal(data, &bookJSON)
	assert.Equal(t, testBook, bookJSON["data"], "body should match")
}

func TestBookHandler_GetBook_Not_Found(t *testing.T) {
	ctx := context.Background()
	handler := getBookHandler()

	mockService := NewMockBookService(t)
	testBookID := int64(1)
	mockService.EXPECT().GetBookByID(ctx, testBookID).Return(book.Book{}, book.ErrNotFound)
	injectBookMocks(handler, mockService)

	request := httptest.NewRequest("GET", "/v1/books/1", nil)
	request.SetPathValue("bookID", strconv.Itoa(int(testBookID)))
	recorder := httptest.NewRecorder()
	err := handler.GetBook(ctx, recorder, request)
	require.Error(t, err, "should get not found error")
	assert.ErrorIs(t, err, apiErrors.ErrNotFound, "should not found book by id")
}

func TestBookHandler_GetBook_InvalidBookID(t *testing.T) {
	ctx := context.Background()
	handler := getBookHandler()

	request := httptest.NewRequest("GET", "/v1/books/one", nil)
	request.SetPathValue("bookID", "one")
	recorder := httptest.NewRecorder()
	err := handler.GetBook(ctx, recorder, request)
	require.Error(t, err, "should not get a book")
	assert.ErrorAs(t, err, &apiErrors.ValidationError{}, "should get a validation error")
}

func TestBookHandler_GetBook_ServiceError(t *testing.T) {
	ctx := context.Background()
	handler := getBookHandler()

	mockService := NewMockBookService(t)
	testBookID := int64(1)
	serviceError := errors.New("service error")
	mockService.EXPECT().GetBookByID(ctx, testBookID).Return(book.Book{}, serviceError)
	injectBookMocks(handler, mockService)

	request := httptest.NewRequest("GET", "/v1/books/1", nil)
	request.SetPathValue("bookID", strconv.Itoa(int(testBookID)))
	recorder := httptest.NewRecorder()
	err := handler.GetBook(ctx, recorder, request)
	require.Error(t, err, "should not get a book")
	assert.ErrorIs(t, err, serviceError, "should get service error")
}

func TestBookHandler_GetBooks(t *testing.T) {
	ctx := context.Background()
	handler := getBookHandler()

	pageNumber := "1"
	pageSize := "100"
	totalItems := int64(10)
	pageNumberNum, _ := strconv.Atoi(pageNumber)
	values := map[string][]string{"page": {pageNumber}, "size": {pageSize}}
	pageRequest, _ := paging.NewPageRequest(values)
	lookupItem := getTestLookupItem()
	page := paging.NewPage(pageRequest, totalItems, []book.LookupItem{lookupItem})

	mockService := NewMockBookService(t)
	mockService.EXPECT().GetBooks(ctx, mock.Anything, mock.Anything, mock.Anything).Return(page, nil)
	injectBookMocks(handler, mockService)

	request := httptest.NewRequest("GET", "/v1/books?page=1&size=10&sort=id,ASC&tag=1&author=1", nil)
	recorder := httptest.NewRecorder()
	err := handler.GetBooks(ctx, recorder, request)
	require.NoError(t, err, "should get a page of books")

	result := recorder.Result()
	defer result.Body.Close()
	require.Equal(t, http.StatusOK, result.StatusCode, "should get a 200 OK response")

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err, "should read body")
	var bookPage map[string]paging.Page[book.LookupItem]
	_ = json.Unmarshal(data, &bookPage)
	pageData := bookPage["data"]
	assert.Equal(t, int64(pageNumberNum), pageData.Page, "page number should match")
	assert.Equal(t, int64(1), pageData.Size, "page size should match")
	assert.Len(t, pageData.Content, int(pageData.Size), "page content size should match page size")
	assert.Equal(t, int64(1), pageData.TotalPages, "total pages should match")
	assert.Equal(t, totalItems, pageData.TotalItems, "total items should match")
	assert.Equal(t, lookupItem, pageData.Content[0], "lookup item content should match")
}

func TestBookHandler_GetBooks_ServiceError(t *testing.T) {
	ctx := context.Background()
	handler := getBookHandler()

	response := paging.Page[book.LookupItem]{}
	serviceError := errors.New("service error")
	mockService := NewMockBookService(t)
	mockService.EXPECT().GetBooks(ctx, mock.Anything, mock.Anything, mock.Anything).Return(response, serviceError)
	injectBookMocks(handler, mockService)

	request := httptest.NewRequest("GET", "/v1/books?page=1&size=10", nil)
	recorder := httptest.NewRecorder()
	err := handler.GetBooks(ctx, recorder, request)
	require.Error(t, err, "should not get books")
	assert.ErrorIs(t, err, serviceError, "should get service error")
}

func TestBookHandler_GetBooks_PageRequestError(t *testing.T) {
	ctx := context.Background()
	handler := getBookHandler()

	request := httptest.NewRequest("GET", "/v1/books?page=one&size=ten", nil)
	recorder := httptest.NewRecorder()
	err := handler.GetBooks(ctx, recorder, request)
	require.Error(t, err, "should get a page of books")
	assert.ErrorAs(t, err, &apiErrors.ValidationError{}, "should get a validation error")
}

func TestBookHandler_GetBooks_SortError(t *testing.T) {
	ctx := context.Background()
	handler := getBookHandler()

	request := httptest.NewRequest("GET", "/v1/books?page=1&size=10&sort=book_wheels,ASC", nil)
	recorder := httptest.NewRecorder()
	err := handler.GetBooks(ctx, recorder, request)
	require.Error(t, err, "should get a page of books")
	assert.ErrorAs(t, err, &apiErrors.ValidationError{}, "should get a validation error")
}

func TestBookHandler_GetBooks_FilterError(t *testing.T) {
	ctx := context.Background()
	handler := getBookHandler()

	request := httptest.NewRequest("GET", "/v1/books?page=1&size=10&sort=id,ASC&tag=one", nil)
	recorder := httptest.NewRecorder()
	err := handler.GetBooks(ctx, recorder, request)
	require.Error(t, err, "should get a page of books")
	assert.ErrorAs(t, err, &apiErrors.ValidationError{}, "should get a validation error")
}

func getBookHandler() *BookHandler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return NewBookHandler(logger, nil)
}

func injectBookMocks(service *BookHandler, bookService *MockBookService) {
	service.bookService = bookService
}

func getTestBook() book.Book {
	bookPubDate, _ := time.Parse(time.DateOnly, bookPubDate)
	createdAt := time.Date(2025, time.April, 10, 9, 15, 10, 0, time.UTC)
	updatedAt := time.Date(2025, time.April, 15, 10, 25, 15, 0, time.UTC)
	return book.Book{
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

func getTestLookupItem() book.LookupItem {
	bookPubDate, _ := time.Parse(time.DateOnly, bookPubDate)
	return book.LookupItem{
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
