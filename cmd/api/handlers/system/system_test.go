package system

import (
	"context"
	"errors"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestController_RegisterRoutes(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))
	testRegistrar := handlers.RouteRegistrarMock{}
	h := NewController(logger)
	h.RegisterRoutes(&testRegistrar)

	assert.True(t, testRegistrar.IsRouteRegistered("GET /livez", h.LivenessProbe))
	assert.True(t, testRegistrar.IsRouteRegistered("GET /readyz", h.ReadinessProbe))
}

func TestLivenessProbe(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))
	controller := Controller{logger: logger}

	req := httptest.NewRequest("GET", "/livez", nil)
	w := httptest.NewRecorder()
	err := controller.LivenessProbe(context.Background(), w, req)

	if assert.NoError(t, err, "should pass liveness probe") {
		result := w.Result()
		defer result.Body.Close()
		assert.Equal(t, http.StatusOK, w.Code)
		bytes, err := io.ReadAll(result.Body)
		if assert.NoError(t, err, "should not return error") {
			assert.JSONEq(t, `{"status":"OK"}`, string(bytes))
		}
	}
}

func TestReadinessProbeFailed(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))

	// online and healthy
	dbHealthChecker := testHealthChecker{
		id: "db", delay: 100 * time.Millisecond, alive: true,
	}
	// offline
	minioHealthChecker := testHealthChecker{
		id: "minio", delay: 1 * time.Second, alive: false,
	}

	controller := NewController(logger, &dbHealthChecker, &minioHealthChecker)
	verifyHealthCheck(t, controller, `{"db":"OK", "minio":"FAIL"}`)
}

func TestReadinessProbeSuccess(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))

	// online and healthy
	dbHealthChecker := testHealthChecker{
		id: "db", delay: 100 * time.Millisecond, alive: true,
	}
	// online and healthy
	minioHealthChecker := testHealthChecker{
		id: "minio", delay: 1 * time.Second, alive: true,
	}

	controller := NewController(logger, &dbHealthChecker, &minioHealthChecker)
	verifyHealthCheck(t, controller, `{"db":"OK", "minio":"OK"}`)
}

func verifyHealthCheck(t *testing.T, controller *Controller, expectedResponse string) {
	req := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()
	err := controller.ReadinessProbe(context.Background(), w, req)
	if assert.NoError(t, err, "should return readiness probe result") {
		result := w.Result()
		defer result.Body.Close()
		if strings.Contains(expectedResponse, (string)(statusFail)) {
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		} else {
			assert.Equal(t, http.StatusOK, w.Code)
		}
		bytes, err := io.ReadAll(result.Body)
		if assert.NoError(t, err, "should not return error") {
			assert.JSONEq(t, expectedResponse, string(bytes))
		}
	}
}

type testHealthChecker struct {
	id    string
	delay time.Duration
	alive bool
}

func (t *testHealthChecker) HealthCheck(ctx context.Context) error {
	<-time.After(t.delay)
	var err error
	if !t.alive {
		err = errors.New("system offline")
	}

	return err
}

func (t *testHealthChecker) HealthCheckID() string {
	return t.id
}
