package main

import (
	"context"
	"errors"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
)

const (
	headerTestMiddleware = "X-Middleware"
)

func TestRouter_Handle(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	testData := `{"data":"test"}`

	r := NewRouter(logger, nil, config.HTTPConfig{})
	clear(r.mw) // disable all application-wide middlewares
	r.RegisterRoute(http.MethodGet, "/v1", "/group-test", getTestHandlerNoError(testData))
	r.RegisterRoute(http.MethodGet, "", "/no-group-test", getTestHandlerNoError(testData))
	r.RegisterRoute(http.MethodGet, "", "/error", getTestHandlerError())

	svr := httptest.NewServer(r.GetHandler())
	defer svr.Close()
	checkResultNoError(t, svr.Client(), svr.URL+"/v1/group-test", testData)
	checkResultNoError(t, svr.Client(), svr.URL+"/no-group-test", testData)
	checkResultError(t, svr.Client(), svr.URL+"/error")

	var manuallyRegisteredRoutesCount int32 = 3
	assert.GreaterOrEqual(t, r.routesCount.Load(), manuallyRegisteredRoutesCount)
}

func TestRouter_MiddlewareRegistration(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	testData := `{"data":"test"}`
	applicationMiddleware, applicationMiddlewareCallsCount := getMockMiddleware("applicationWideMiddleware")
	handlerMiddleware, handlerMiddlewareCallsCount := getMockMiddleware("handlerSpecificMiddleware")

	r := NewRouter(logger, nil, config.HTTPConfig{})
	r.AddApplicationMiddleware(applicationMiddleware)
	r.RegisterRoute(http.MethodGet, "", "/no-handler-middleware", getTestHandlerNoError(testData))
	r.RegisterRoute(http.MethodGet, "", "/handler-middleware", getTestHandlerNoError(testData), handlerMiddleware)

	svr := httptest.NewServer(r.GetHandler())
	defer svr.Close()
	checkResultNoError(t, svr.Client(), svr.URL+"/no-handler-middleware", testData)
	checkResultNoError(t, svr.Client(), svr.URL+"/handler-middleware", testData)

	var manuallyRegisteredRoutesCount int32 = 2
	var expectedApplicationMiddlewareCallsCount int32 = 2 // one call for each endpoint
	var expectedHandlerMiddlewareCallsCount int32 = 1     // one call for one endpoint
	require.GreaterOrEqual(t, r.routesCount.Load(), manuallyRegisteredRoutesCount)
	require.Equal(t, expectedApplicationMiddlewareCallsCount, applicationMiddlewareCallsCount.Load())
	require.Equal(t, expectedHandlerMiddlewareCallsCount, handlerMiddlewareCallsCount.Load())
}

func TestRouter_MiddlewareExecutionOrder(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	middleware01Name := "01"
	middleware02Name := "02"
	applicationMiddleware01, _ := getMockMiddleware(middleware01Name)
	applicationMiddleware02, _ := getMockMiddleware(middleware02Name)

	r := NewRouter(logger, nil, config.HTTPConfig{})
	r.AddApplicationMiddleware(applicationMiddleware02)
	r.AddApplicationMiddleware(applicationMiddleware01)
	r.RegisterRoute(http.MethodGet, "", "/middleware", getTestHandlerNoError(`{"data":"test"}`))

	svr := httptest.NewServer(r.GetHandler())
	defer svr.Close()

	resp, err := svr.Client().Get(svr.URL + "/middleware")
	if assert.NotNil(t, resp.Body) {
		defer resp.Body.Close()
	}
	require.NoError(t, err, "should return response")
	headerValues := resp.Header.Values(headerTestMiddleware)
	require.Len(t, headerValues, 2)
	require.Equal(t, middleware02Name, headerValues[0])
	require.Equal(t, middleware01Name, headerValues[1])
}

func getMockMiddleware(name string) (handlers.Middleware, *atomic.Int32) {
	var callsCount atomic.Int32
	return func(next handlers.HTTPHandler) handlers.HTTPHandler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			callsCount.Add(1)

			w.Header().Add(headerTestMiddleware, name)
			return next(ctx, w, r)
		}
	}, &callsCount
}

func getTestHandlerNoError(testData string) handlers.HTTPHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testData))
		return err
	}
}

func getTestHandlerError() handlers.HTTPHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return errors.New("internal error")
	}
}

func checkResultNoError(t *testing.T, client *http.Client, URL string, expectedResult string) {
	resp, err := client.Get(URL)
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if assert.NoError(t, err, "should return response") {
		responseData, err2 := io.ReadAll(resp.Body)
		if assert.NoError(t, err2, "body reading error") {
			assert.Equal(t, expectedResult, string(responseData))
		}
	}
}

func checkResultError(t *testing.T, client *http.Client, URL string) {
	resp, err := client.Get(URL)
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if assert.NoError(t, err, "should return response") {
		response, err2 := io.ReadAll(resp.Body)
		if assert.NoError(t, err2, "body reading error") {
			assert.Equal(t, http.StatusText(http.StatusInternalServerError)+"\n", string(response))
		}
	}
}
