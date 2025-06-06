package v1

import (
	"context"
	"encoding/json"
	"errors"
	apiErrors "github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/domain/publisher"
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
	publisherID   = int64(1)
	publisherName = "OReilly"
)

func TestPublisherController_RegisterRoutes(t *testing.T) {
	testRegistrar := handlers.RouteRegistrarMock{}
	cnt := getPublisherController()
	cnt.RegisterRoutes(&testRegistrar)

	assert.True(t, testRegistrar.IsRouteRegistered("GET /v1/publishers", cnt.GetPublishers))
}

func TestPublisherController_GetPublishers(t *testing.T) {
	ctx := context.Background()
	controller := getPublisherController()

	pageNumber := "1"
	pageSize := "100"
	totalItems := int64(10)
	pageNumberNum, _ := strconv.Atoi(pageNumber)
	values := map[string][]string{"page": {pageNumber}, "size": {pageSize}}
	pageRequest, _ := paging.NewPageRequest(values)
	lookupItem := getTestPublisherLookupItem()
	page := paging.NewPage(pageRequest, totalItems, []publisher.LookupItem{getTestPublisherLookupItem()})

	mockService := NewMockPublisherService(t)
	mockService.EXPECT().GetPublishers(ctx, mock.Anything, mock.Anything).Return(page, nil)
	injectPublisherMocks(controller, mockService)

	request := httptest.NewRequest("GET", "/v1/publishers?page=1&size=10&sort=id,ASC", nil)
	recorder := httptest.NewRecorder()
	err := controller.GetPublishers(ctx, recorder, request)
	require.NoError(t, err, "should get a page of publishers")

	result := recorder.Result()
	defer result.Body.Close()
	require.Equal(t, http.StatusOK, result.StatusCode, "should get a 200 OK response")

	data, err := io.ReadAll(result.Body)
	require.NoError(t, err, "should read body")
	var PublisherPage map[string]paging.Page[publisher.LookupItem]
	_ = json.Unmarshal(data, &PublisherPage)
	pageData := PublisherPage["data"]
	assert.Equal(t, int64(pageNumberNum), pageData.Page, "page number should match")
	assert.Equal(t, int64(1), pageData.Size, "page size should match")
	assert.Len(t, pageData.Content, int(pageData.Size), "page content size should match page size")
	assert.Equal(t, int64(1), pageData.TotalPages, "total pages should match")
	assert.Equal(t, totalItems, pageData.TotalItems, "total items should match")
	assert.Equal(t, lookupItem, pageData.Content[0], "lookup item content should match")
}

func TestPublisherController_GetPublishers_ServiceError(t *testing.T) {
	ctx := context.Background()
	controller := getPublisherController()

	response := paging.Page[publisher.LookupItem]{}
	serviceError := errors.New("service error")
	mockService := NewMockPublisherService(t)
	mockService.EXPECT().GetPublishers(ctx, mock.Anything, mock.Anything).Return(response, serviceError)
	injectPublisherMocks(controller, mockService)

	request := httptest.NewRequest("GET", "/v1/publishers?page=1&size=10", nil)
	recorder := httptest.NewRecorder()
	err := controller.GetPublishers(ctx, recorder, request)
	require.Error(t, err, "should not get publishers")
	require.ErrorIs(t, err, serviceError, "should get service error")
}

func TestPublisherController_GetPublishers_PageRequestError(t *testing.T) {
	ctx := context.Background()
	controller := getPublisherController()

	request := httptest.NewRequest("GET", "/v1/publishers?page=one&size=ten", nil)
	recorder := httptest.NewRecorder()
	err := controller.GetPublishers(ctx, recorder, request)
	require.Error(t, err)
	require.ErrorAs(t, err, &apiErrors.ValidationError{}, "should get a validation error")
}

func TestPublisherController_GetPublishers_SortError(t *testing.T) {
	ctx := context.Background()
	controller := getPublisherController()

	request := httptest.NewRequest("GET", "/v1/publishers?page=1&size=10&sort=age,ASC", nil)
	recorder := httptest.NewRecorder()
	err := controller.GetPublishers(ctx, recorder, request)
	require.Error(t, err, "should get a page of publishers")
	assert.ErrorAs(t, err, &apiErrors.ValidationError{}, "should get a validation error")
}

func getPublisherController() *PublisherController {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return NewPublisherController(logger, nil)
}

func injectPublisherMocks(controller *PublisherController, PublisherService *MockPublisherService) {
	controller.service = PublisherService
}

func getTestPublisherLookupItem() publisher.LookupItem {
	return publisher.LookupItem{
		ID:   publisherID,
		Name: publisherName,
	}
}
