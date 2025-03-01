package filetype

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sdreger/lib-manager-go/internal/paging"
)

type DBStore struct {
	db *sqlx.DB
}

func NewDBStore(db *sqlx.DB) *DBStore {
	return &DBStore{db: db}
}

func (s *DBStore) Lookup(ctx context.Context, page paging.PageRequest, sort paging.Sort) ([]LookupItem, int64, error) {

	query := fmt.Sprintf("SELECT id, name FROM ebook.file_types ORDER BY %s LIMIT $1 OFFSET $2",
		sort.GetOrderBy())

	var rows []lookupEntity
	err := s.db.SelectContext(ctx, &rows, query, page.Limit(), page.Offset())
	if err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	err = s.db.GetContext(ctx, &total, "SELECT count(id) FROM ebook.file_types")
	if err != nil {
		return nil, 0, err
	}

	items := make([]LookupItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, LookupItem(row))
	}

	return items, total, nil
}
