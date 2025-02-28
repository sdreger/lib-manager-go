package book

import (
	"context"
	"errors"
	"github.com/sdreger/lib-manager-go/internal/paging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"strconv"
	"testing"
)

func TestService_GetById(t *testing.T) {
	ctx := context.Background()
	service := getService()

	mockStore := NewMockStore(t)
	mockStore.EXPECT().GetByID(ctx, bookID).Return(getTestBook(), nil).Once()
	injectMocks(service, mockStore)

	book, err := service.GetBookByID(ctx, bookID)
	if assert.NoError(t, err, "should get book by id") {
		assert.Equal(t, getTestBook(), book, "books should be equal")
	}
}

func TestService_GetById_NotFound(t *testing.T) {
	ctx := context.Background()
	service := getService()

	mockStore := NewMockStore(t)
	mockStore.EXPECT().GetByID(ctx, bookID).Return(Book{}, ErrNotFound).Once()
	injectMocks(service, mockStore)

	_, err := service.GetBookByID(ctx, bookID)
	if assert.Error(t, err, "should not found book by id") {
		assert.ErrorIs(t, err, ErrNotFound, "should not found book by id")
	}
}

func TestService_GetBooks_Success(t *testing.T) {
	ctx := context.Background()
	service := getService()

	pageNumber := "1"
	pageSize := "1"
	pageSizeNum, _ := strconv.Atoi(pageSize)
	values := map[string][]string{"page": {pageNumber}, "size": {pageSize}}
	pageRequest, _ := paging.NewPageRequest(values)
	sort, _ := paging.NewSort(values, AllowedSortFields)
	filter, _ := NewFilter(values)

	mockStore := NewMockStore(t)
	totalItems := int64(5)
	lookupItem := getTestLookupItem()
	response := []LookupItem{lookupItem}
	mockStore.EXPECT().Lookup(ctx, pageRequest, sort, filter).Return(response, totalItems, nil).Once()
	injectMocks(service, mockStore)

	page, err := service.GetBooks(ctx, pageRequest, sort, filter)
	if assert.NoError(t, err, "should find books") {
		content := page.Content
		assert.Len(t, content, 1)
		book := content[0]
		assert.Equal(t, lookupItem, book, "books should be equal")
		assert.Equal(t, int64(pageSizeNum), page.Page)
		assert.Len(t, content, int(page.Size))
		assert.Equal(t, totalItems/int64(pageSizeNum), page.TotalPages)
		assert.Equal(t, totalItems, page.TotalItems)
	}
}

func TestService_GetBooks_Failure(t *testing.T) {
	ctx := context.Background()
	service := getService()

	pageNumber := "1"
	pageSize := "1"
	values := map[string][]string{"page": {pageNumber}, "size": {pageSize}}
	pageRequest, _ := paging.NewPageRequest(values)
	sort, _ := paging.NewSort(values, AllowedSortFields)
	filter, _ := NewFilter(values)

	mockStore := NewMockStore(t)
	storeError := errors.New("some error")
	mockStore.EXPECT().Lookup(ctx, pageRequest, sort, filter).Return(nil, 0, storeError).Once()
	injectMocks(service, mockStore)

	page, err := service.GetBooks(ctx, pageRequest, sort, filter)
	require.Error(t, err, "should get an error")
	require.ErrorIs(t, err, storeError, "should get the correct error")
	assert.Empty(t, page)
}

func getService() *Service {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return NewService(logger, nil)
}

func injectMocks(service *Service, store *MockStore) {
	service.store = store
}
