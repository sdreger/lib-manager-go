package book

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Store interface {
	GetByID(ctx context.Context, bookID int64) (Book, error)
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
