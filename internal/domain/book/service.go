package book

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sdreger/lib-manager-go/internal/paging"
	"log/slog"
)

type Store interface {
	GetByID(ctx context.Context, bookID int64) (Book, error)
	Lookup(
		ctx context.Context,
		page paging.PageRequest,
		sort paging.Sort,
		filter Filter,
	) ([]LookupItem, int64, error)
}

type Service struct {
	logger *slog.Logger
	store  Store
}

func NewService(logger *slog.Logger, db *sqlx.DB) *Service {
	return &Service{
		logger: logger,
		store:  NewDBStore(db),
	}
}

// GetBookByID - returns a book from the database if it exists
func (s Service) GetBookByID(ctx context.Context, bookID int64) (Book, error) {
	return s.store.GetByID(ctx, bookID)
}

// GetBooks - returns a requested page of books based on provided filter values
func (s Service) GetBooks(ctx context.Context, pageRequest paging.PageRequest, sort paging.Sort, filter Filter) (
	paging.Page[LookupItem], error) {

	lookupItems, totalElements, err := s.store.Lookup(ctx, pageRequest, sort, filter)
	if err != nil {
		return paging.Page[LookupItem]{}, err
	}

	return paging.NewPage(pageRequest, totalElements, lookupItems), nil
}
