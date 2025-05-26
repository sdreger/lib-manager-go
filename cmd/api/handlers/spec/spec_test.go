package spec

import (
	"context"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestController_RegisterRoutes(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))
	testRegistrar := handlers.RouteRegistrarMock{}
	h := NewController(logger)
	h.RegisterRoutes(&testRegistrar)

	assert.True(t, testRegistrar.IsRouteRegistered("GET "+swaggerURLPrefix, h.SwaggerUI))
	assert.True(t, testRegistrar.IsRouteRegistered("GET "+specURLPrefix, h.APISpec))
}

func TestSwaggerUI(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))
	controller := Controller{logger: logger}

	req := httptest.NewRequest("GET", swaggerURLPrefix, nil)
	w := httptest.NewRecorder()
	err := controller.SwaggerUI(context.Background(), w, req)
	if assert.NoError(t, err, "should get SwaggerUI index page") {
		result := w.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusOK, w.Code)
	}

	req = httptest.NewRequest("GET", swaggerURLPrefix+"swagger-ui-bundle.js", nil)
	w = httptest.NewRecorder()
	err = controller.SwaggerUI(context.Background(), w, req)
	if assert.NoError(t, err, "should get SwaggerUI JS bundle") {
		result := w.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestOpenAPISpec(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))
	controller := Controller{logger: logger}

	specFile, err := apiSpec.Open(specFileName)
	require.NoError(t, err)
	specFileContent, err := io.ReadAll(specFile)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", specURLPrefix+specFileName, nil)
	w := httptest.NewRecorder()
	err = controller.APISpec(context.Background(), w, req)
	if assert.NoError(t, err, "should get OpenAPI spec file") {
		result := w.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusOK, w.Code)
		bytes, err := io.ReadAll(result.Body)
		if assert.NoError(t, err, "should not return error") {
			assert.YAMLEq(t, string(specFileContent), string(bytes))
		}
	}
}
