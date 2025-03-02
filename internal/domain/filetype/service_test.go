package filetype

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

const (
	fileTypeID   = int64(1)
	fileTypeName = "pdf"
)

func TestService_GetFileTypes_Success(t *testing.T) {
	ctx := context.Background()
	service := getService()

	pageNumber := 1
	pageSize := 1
	values := map[string][]string{
		"page": {strconv.Itoa(pageNumber)},
		"size": {strconv.Itoa(pageSize)},
		"sort": {"name,asc"},
	}
	pageRequest, _ := paging.NewPageRequest(values)
	sort, _ := paging.NewSort(values, AllowedSortFields)

	mockStore := NewMockStore(t)
	totalItems := int64(5)
	lookupItem := getTestLookupItem()
	response := []LookupItem{lookupItem}
	mockStore.EXPECT().Lookup(ctx, pageRequest, sort).Return(response, totalItems, nil).Once()
	injectMocks(service, mockStore)

	page, err := service.GetFileTypes(ctx, pageRequest, sort)
	if assert.NoError(t, err, "should find file types") {
		content := page.Content
		assert.Len(t, content, 1)
		fileType := content[0]
		assert.Equal(t, lookupItem, fileType, "file types should be equal")
		assert.Equal(t, int64(pageSize), page.Page)
		assert.Len(t, content, int(page.Size))
		assert.Equal(t, totalItems/int64(pageSize), page.TotalPages)
		assert.Equal(t, totalItems, page.TotalItems)
	}
}

func TestService_GetFileTypes_Failure(t *testing.T) {
	ctx := context.Background()
	service := getService()

	pageNumber := "1"
	pageSize := "1"
	values := map[string][]string{"page": {pageNumber}, "size": {pageSize}, "sort": {"name,asc"}}
	pageRequest, _ := paging.NewPageRequest(values)
	sort, _ := paging.NewSort(values, AllowedSortFields)

	mockStore := NewMockStore(t)
	storeError := errors.New("some error")
	mockStore.EXPECT().Lookup(ctx, pageRequest, sort).Return(nil, 0, storeError).Once()
	injectMocks(service, mockStore)

	page, err := service.GetFileTypes(ctx, pageRequest, sort)
	require.Error(t, err, "should get an error")
	require.ErrorIs(t, err, storeError, "should get the correct error")
	assert.Empty(t, page)
}

func getTestLookupItem() LookupItem {
	return LookupItem{
		ID:   fileTypeID,
		Name: fileTypeName,
	}
}

func getService() *Service {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return NewService(logger, nil)
}

func injectMocks(service *Service, store *MockStore) {
	service.store = store
}
