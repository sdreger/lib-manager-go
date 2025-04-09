package paging

import (
	"fmt"
	"github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSort(t *testing.T) {
	tt := []struct {
		sortString      string
		expectedOrderBy string
		err             bool
	}{
		{sortString: "", expectedOrderBy: "id ASC", err: false},
		{sortString: "id,asc", expectedOrderBy: "id ASC", err: false},
		{sortString: "id,desc", expectedOrderBy: "id DESC", err: false},
		{sortString: "id,ASC", expectedOrderBy: "id ASC", err: false},
		{sortString: "id,DESC", expectedOrderBy: "id DESC", err: false},
		{sortString: "title,asc", expectedOrderBy: "title ASC", err: false},
		{sortString: "title,desc", expectedOrderBy: "title DESC", err: false},
		{sortString: "subtitle,asc", expectedOrderBy: "subtitle ASC", err: false},
		{sortString: "subtitle,desc", expectedOrderBy: "subtitle DESC", err: false},
		{sortString: "isbn10,asc", expectedOrderBy: "isbn10 ASC", err: false},
		{sortString: "isbn10,desc", expectedOrderBy: "isbn10 DESC", err: false},
		{sortString: "isbn13,asc", expectedOrderBy: "isbn13 ASC", err: false},
		{sortString: "isbn13,desc", expectedOrderBy: "isbn13 DESC", err: false},
		{sortString: "asin,asc", expectedOrderBy: "asin ASC", err: false},
		{sortString: "asin,desc", expectedOrderBy: "asin DESC", err: false},
		{sortString: "pages,asc", expectedOrderBy: "pages ASC", err: false},
		{sortString: "pages,desc", expectedOrderBy: "pages DESC", err: false},
		{sortString: "edition,asc", expectedOrderBy: "edition ASC", err: false},
		{sortString: "edition,desc", expectedOrderBy: "edition DESC", err: false},
		{sortString: "pub_date,asc", expectedOrderBy: "pub_date ASC", err: false},
		{sortString: "pub_date,desc", expectedOrderBy: "pub_date DESC", err: false},
		{sortString: "book_file_size,asc", expectedOrderBy: "book_file_size ASC", err: false},
		{sortString: "book_file_size,desc", expectedOrderBy: "book_file_size DESC", err: false},
		{sortString: "created_at,asc", expectedOrderBy: "created_at ASC", err: false},
		{sortString: "created_at,desc", expectedOrderBy: "created_at DESC", err: false},
		{sortString: "updated_at,asc", expectedOrderBy: "updated_at ASC", err: false},
		{sortString: "updated_at,desc", expectedOrderBy: "updated_at DESC", err: false},
		{sortString: "description,asc", expectedOrderBy: "", err: true},
		{sortString: "id,descending", expectedOrderBy: "", err: true},
		{sortString: "id,ascending,descending", expectedOrderBy: "", err: true},
		{sortString: "id,desc; DROP table books;", expectedOrderBy: "", err: true},
	}

	allowedSortFields := []string{
		"id", "title", "subtitle", "isbn10", "isbn13", "asin", "pages", "edition",
		"pub_date", "book_file_size", "created_at", "updated_at",
	}
	for _, tc := range tt {
		t.Run(tc.sortString, func(t *testing.T) {
			values := map[string][]string{"sort": {tc.sortString}}
			sort, err := NewSort(values, allowedSortFields)
			if tc.err {
				require.Error(t, err)
				assert.ErrorAs(t, err, &errors.ValidationError{})
			} else {
				orderFieldPrefix := "ebook.books"
				require.NoError(t, err)
				expectedOrderBy := fmt.Sprintf("%s.%s", orderFieldPrefix, tc.expectedOrderBy)
				assert.Equal(t, expectedOrderBy, sort.GetOrderBy(orderFieldPrefix))
			}
		})
	}
}
