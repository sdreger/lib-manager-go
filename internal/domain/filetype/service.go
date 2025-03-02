package filetype

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sdreger/lib-manager-go/internal/paging"
	"log/slog"
)

type Store interface {
	Lookup(ctx context.Context, page paging.PageRequest, sort paging.Sort) ([]LookupItem, int64, error)
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

// GetFileTypes - returns a requested page of fileTypes
func (s Service) GetFileTypes(ctx context.Context, pageRequest paging.PageRequest, sort paging.Sort) (
	paging.Page[LookupItem], error) {

	lookupItems, totalElements, err := s.store.Lookup(ctx, pageRequest, sort)
	if err != nil {
		return paging.Page[LookupItem]{}, err
	}

	return paging.NewPage(pageRequest, totalElements, lookupItems), nil
}
