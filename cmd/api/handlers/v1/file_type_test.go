package v1

import (
	"context"
	"encoding/json"
	"errors"
	apiErrors "github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/domain/filetype"
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
)

const (
	fileTypeID   = int64(1)
	fileTypeName = "pdf"
)

func TestFileTypeController_RegisterRoutes(t *testing.T) {
	testRegistrar := handlers.RouteRegistrarMock{}
	cnt := getFileTypeController()
	cnt.RegisterRoutes(&testRegistrar)

	assert.True(t, testRegistrar.IsRouteRegistered("GET /v1/file_types", cnt.GetFileTypes))
}

func TestFileTypeController_GetFileTypes(t *testing.T) {
	ctx := context.Background()
	controller := getFileTypeController()

	pageNumber := "1"
	pageSize := "100"
	totalItems := int64(10)
	pageNumberNum, _ := strconv.Atoi(pageNumber)
	values := map[string][]string{"page": {pageNumber}, "size": {pageSize}}
	pageRequest, _ := paging.NewPageRequest(values)
	lookupItem := getTestFileTypeLookupItem()
	page := paging.NewPage(pageRequest, totalItems, []filetype.LookupItem{getTestFileTypeLookupItem()})

	mockService := NewMockFileTypeService(t)
	mockService.EXPECT().GetFileTypes(ctx, mock.Anything, mock.Anything).Return(page, nil)
	injectFileTypeMocks(controller, mockService)

	request := httptest.NewRequest("GET", "/v1/file_types?page=1&size=10&sort=id,ASC", nil)
	recorder := httptest.NewRecorder()
	err := controller.GetFileTypes(ctx, recorder, request)
	require.NoError(t, err, "should get a page of file types")

	result := recorder.Result()
	defer result.Body.Close()
	require.Equal(t, http.StatusOK, result.StatusCode, "should get a 200 OK response")

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err, "should read body")
	var FileTypePage map[string]paging.Page[filetype.LookupItem]
	_ = json.Unmarshal(data, &FileTypePage)
	pageData := FileTypePage["data"]
	assert.Equal(t, int64(pageNumberNum), pageData.Page, "page number should match")
	assert.Equal(t, int64(1), pageData.Size, "page size should match")
	assert.Len(t, pageData.Content, int(pageData.Size), "page content size should match page size")
	assert.Equal(t, int64(1), pageData.TotalPages, "total pages should match")
	assert.Equal(t, totalItems, pageData.TotalItems, "total items should match")
	assert.Equal(t, lookupItem, pageData.Content[0], "lookup item content should match")
}

func TestFileTypeController_GetFileTypes_ServiceError(t *testing.T) {
	ctx := context.Background()
	controller := getFileTypeController()

	response := paging.Page[filetype.LookupItem]{}
	serviceError := errors.New("service error")
	mockService := NewMockFileTypeService(t)
	mockService.EXPECT().GetFileTypes(ctx, mock.Anything, mock.Anything).Return(response, serviceError)
	injectFileTypeMocks(controller, mockService)

	request := httptest.NewRequest("GET", "/v1/file_types?page=1&size=10", nil)
	recorder := httptest.NewRecorder()
	err := controller.GetFileTypes(ctx, recorder, request)
	require.Error(t, err, "should not get file types")
	require.ErrorIs(t, err, serviceError, "should get service error")
}

func TestFileTypeController_GetFileTypes_PageRequestError(t *testing.T) {
	ctx := context.Background()
	controller := getFileTypeController()

	request := httptest.NewRequest("GET", "/v1/file_types?page=one&size=ten", nil)
	recorder := httptest.NewRecorder()
	err := controller.GetFileTypes(ctx, recorder, request)
	require.Error(t, err)
	require.ErrorAs(t, err, &apiErrors.ValidationError{}, "should get a validation error")
}

func TestFileTypeController_GetFileTypes_SortError(t *testing.T) {
	ctx := context.Background()
	controller := getFileTypeController()

	request := httptest.NewRequest("GET", "/v1/file_types?page=1&size=10&sort=age,ASC", nil)
	recorder := httptest.NewRecorder()
	err := controller.GetFileTypes(ctx, recorder, request)
	require.Error(t, err, "should get a page of file types")
	assert.ErrorAs(t, err, &apiErrors.ValidationError{}, "should get a validation error")
}

func getFileTypeController() *FileTypeController {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return NewFileTypeController(logger, nil)
}

func injectFileTypeMocks(controller *FileTypeController, FileTypeService *MockFileTypeService) {
	controller.service = FileTypeService
}

func getTestFileTypeLookupItem() filetype.LookupItem {
	return filetype.LookupItem{
		ID:   fileTypeID,
		Name: fileTypeName,
	}
}
