package system

import (
	"context"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/stretchr/testify/assert"
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

	assert.True(t, testRegistrar.IsRouteRegistered("GET /health", h.HealthProbe))
}

func TestHealthProbe(t *testing.T) {
	controller := Controller{logger: slog.New(slog.NewJSONHandler(os.Stdout, nil))}

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	err := controller.HealthProbe(context.Background(), w, req)

	if assert.NoError(t, err, "should pass health probe") {
		result := w.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusOK, w.Code)
		bytes, err := io.ReadAll(result.Body)
		if assert.NoError(t, err, "should not return error") {
			assert.JSONEq(t, `{"status":"OK"}`, string(bytes))
		}
	}
}
