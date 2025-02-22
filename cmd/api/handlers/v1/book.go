package v1

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	apiErrors "github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	book "github.com/sdreger/lib-manager-go/internal/domain/book"
	"github.com/sdreger/lib-manager-go/internal/paging"
	"github.com/sdreger/lib-manager-go/internal/response"
	"log/slog"
	"net/http"
	"strconv"
)

type BookService interface {
	GetBookByID(ctx context.Context, bookID int64) (book.Book, error)
	GetBooks(
		ctx context.Context,
		pageRequest paging.PageRequest,
		sort paging.Sort,
		filter book.Filter,
	) (paging.Page[book.LookupItem], error)
}

type BookHandler struct {
	logger      *slog.Logger
	bookService BookService
}

func NewBookHandler(logger *slog.Logger, db *sqlx.DB) *BookHandler {
	return &BookHandler{logger: logger, bookService: book.NewService(logger, db)}
}

func (bh *BookHandler) RegisterHandler(registrar handlers.RouteRegistrar) {
	group := "/v1"
	registrar.RegisterRoute(http.MethodGet, group, "/books", bh.GetBooks)
	registrar.RegisterRoute(http.MethodGet, group, "/books/{bookID}", bh.GetBook)
}

func (bh *BookHandler) GetBook(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	idString := r.PathValue("bookID")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return apiErrors.ValidationError{
			Field:   "bookID",
			Message: "the provided bookID should be a number",
		}
	}

	bookEntry, err := bh.bookService.GetBookByID(ctx, int64(idInt))
	if errors.Is(err, book.ErrNotFound) {
		return apiErrors.ErrNotFound
	}
	if err != nil {
		return err
	}

	return response.RenderDataJSON(w, http.StatusOK, bookEntry)
}

func (bh *BookHandler) GetBooks(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, pageErr := paging.NewPageRequest(r.URL.Query())
	if pageErr != nil {
		return pageErr
	}

	sort, sortErr := paging.NewSort(r.URL.Query())
	if sortErr != nil {
		return sortErr
	}

	filter, filterErr := book.NewFilter(r.URL.Query())
	if filterErr != nil {
		return filterErr
	}

	bookPage, err := bh.bookService.GetBooks(ctx, page, sort, filter)
	if err != nil {
		return err
	}

	return response.RenderDataJSON(w, http.StatusOK, bookPage)
}
