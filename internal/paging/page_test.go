package paging

import (
	"github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewPageRequest(t *testing.T) {
	tt := []struct {
		page           string
		size           string
		expectedPage   int64
		expectedSize   int64
		expectedOffset uint64
		expectedLimit  uint64
		err            bool
	}{
		{page: "", size: "", expectedPage: 0, expectedSize: 10, expectedOffset: 0, expectedLimit: 10, err: false},
		{page: "", size: "10", expectedPage: 0, expectedSize: 10, expectedOffset: 0, expectedLimit: 10, err: false},
		{page: "1", size: "", expectedPage: 1, expectedSize: 10, expectedOffset: 10, expectedLimit: 10, err: false},
		{page: "0", size: "10", expectedPage: 0, expectedSize: 10, expectedOffset: 0, expectedLimit: 10, err: false},
		{page: "5", size: "100", expectedPage: 5, expectedSize: 100, expectedOffset: 500, expectedLimit: 100, err: false},
		{page: "500", size: "10", expectedPage: 500, expectedSize: 10, expectedOffset: 5000, expectedLimit: 10, err: false},
		{page: "one", size: "10", expectedPage: 0, expectedSize: 0, expectedOffset: 0, expectedLimit: 0, err: true},
		{page: "0", size: "ten", expectedPage: 0, expectedSize: 0, expectedOffset: 0, expectedLimit: 0, err: true},
		{page: "0", size: "0", expectedPage: 0, expectedSize: 0, expectedOffset: 0, expectedLimit: 0, err: true},
		{page: "0", size: "-1", expectedPage: 0, expectedSize: 0, expectedOffset: 0, expectedLimit: 0, err: true},
		{page: "-1", size: "10", expectedPage: 0, expectedSize: 0, expectedOffset: 0, expectedLimit: 0, err: true},
	}

	for _, tc := range tt {
		t.Run(tc.page+":"+tc.size, func(t *testing.T) {
			pageRequest, err := NewPageRequest(tc.page, tc.size)
			if tc.err {
				require.Error(t, err)
				assert.ErrorAs(t, err, &errors.ValidationError{})
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedOffset, pageRequest.Offset())
				assert.Equal(t, tc.expectedLimit, pageRequest.Limit())
				assert.Equal(t, tc.expectedPage, pageRequest.page)
				assert.Equal(t, tc.expectedSize, pageRequest.size)
			}
		})
	}
}

func TestNewPage(t *testing.T) {
	tt := []struct {
		requestPage        string
		requestSize        string
		totalItems         int64
		contentItems       []interface{}
		expectedPage       int64
		expectedSize       int64
		expectedTotalPages int64
		expectedTotalItems int64
	}{
		{requestPage: "0", requestSize: "10", totalItems: 120, contentItems: make([]interface{}, 10),
			expectedPage: 0, expectedSize: 10, expectedTotalPages: 12, expectedTotalItems: 120}, // first age
		{requestPage: "2", requestSize: "5", totalItems: 13, contentItems: make([]interface{}, 3),
			expectedPage: 2, expectedSize: 3, expectedTotalPages: 3, expectedTotalItems: 13}, // last page
		{requestPage: "9", requestSize: "5", totalItems: 87, contentItems: make([]interface{}, 5),
			expectedPage: 9, expectedSize: 5, expectedTotalPages: 18, expectedTotalItems: 87}, // middle page
	}

	for _, tc := range tt {
		t.Run(t.Name(), func(t *testing.T) {
			pageRequest, err := NewPageRequest(tc.requestPage, tc.requestSize)
			require.NoError(t, err, "error creating requestPage request")
			page := NewPage(pageRequest, tc.totalItems, tc.contentItems)

			assert.Equal(t, tc.expectedPage, page.Page)
			assert.Equal(t, tc.expectedSize, page.Size)
			assert.Equal(t, tc.expectedTotalPages, page.TotalPages)
			assert.Equal(t, tc.expectedTotalItems, page.TotalItems)
			assert.Equal(t, tc.expectedSize, page.Size)

		})
	}
}
