package v1

import (
	"bytes"
	"context"
	"errors"
	apiErrors "github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/domain/cover"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCoverHandler_RegisterCoverHandler(t *testing.T) {

	testRegistrar := handlers.RouteRegistrarMock{}
	h := getCoverHandler()
	h.RegisterHandler(&testRegistrar)

	assert.True(t, testRegistrar.IsRouteRegistered("GET /v1/covers/{publisherName}/{coverFileName}", h.GetBookCover))
}

func TestCoverHandler_GetCover_Success(t *testing.T) {
	ctx := context.Background()
	handler := getCoverHandler()
	filePathExists := "manning/exists.svg"
	existingContent := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>`

	mockService := NewMockCoverService(t)
	mockService.EXPECT().GetBookCover(ctx, filePathExists).Return(bytes.NewBufferString(existingContent), nil)
	injectCoverMocks(handler, mockService)

	request := httptest.NewRequest("GET", "/v1/covers/manning/exists.svg", nil)
	request.SetPathValue("publisherName", "manning")
	request.SetPathValue("coverFileName", "exists.svg")
	recorder := httptest.NewRecorder()
	err := handler.GetBookCover(ctx, recorder, request)
	require.NoError(t, err, "should get a book cover")

	result := recorder.Result()
	defer result.Body.Close()
	require.Equal(t, http.StatusOK, result.StatusCode, "should get a 200 OK response")

	content, err := io.ReadAll(result.Body)
	require.NoError(t, err, "should read body")
	assert.Equal(t, existingContent, string(content))
}

func TestCoverHandler_GetCover_Not_Found(t *testing.T) {
	ctx := context.Background()
	handler := getCoverHandler()
	filePathNotExist := "manning/not-exist.svg"

	mockService := NewMockCoverService(t)
	mockService.EXPECT().GetBookCover(ctx, filePathNotExist).Return(nil, cover.ErrNotFound)
	injectCoverMocks(handler, mockService)

	request := httptest.NewRequest("GET", "/v1/covers/manning/not-exist.svg", nil)
	request.SetPathValue("publisherName", "manning")
	request.SetPathValue("coverFileName", "not-exist.svg")
	recorder := httptest.NewRecorder()
	err := handler.GetBookCover(ctx, recorder, request)
	assert.ErrorIs(t, err, apiErrors.ErrNotFound, "should not found cover")
}

func TestCoverHandler_GetCover_Error(t *testing.T) {
	ctx := context.Background()
	handler := getCoverHandler()

	expectedError := errors.New("some error")
	mockService := NewMockCoverService(t)
	mockService.EXPECT().GetBookCover(ctx, mock.Anything).Return(nil, expectedError)
	injectCoverMocks(handler, mockService)

	request := httptest.NewRequest("GET", "/v1/covers/manning/111111.svg", nil)
	request.SetPathValue("publisherName", "manning")
	request.SetPathValue("coverFileName", "111111.svg")
	recorder := httptest.NewRecorder()
	err := handler.GetBookCover(ctx, recorder, request)
	assert.ErrorIs(t, err, expectedError)
}

func getCoverHandler() *CoverHandler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))
	return NewCoverHandler(logger, nil)
}

func injectCoverMocks(service *CoverHandler, coverService *MockCoverService) {
	service.coverService = coverService
}
