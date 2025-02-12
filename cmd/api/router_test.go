package main

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
	"testing"
)

func TestRouter_Handle(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	testData := `{"data":"test"}`
	errorData := "internal error"

	r := NewRouter(logger, nil)
	r.RegisterRoute(http.MethodGet, "/api/v1", "/group-test", getTestHandlerNoError(testData))
	r.RegisterRoute(http.MethodGet, "", "/no-group-test", getTestHandlerNoError(testData))
	r.RegisterRoute(http.MethodGet, "", "/error", getTestHandlerError(errorData))

	svr := httptest.NewServer(r.GetHandler())
	defer svr.Close()
	checkResultNoError(t, svr.Client(), svr.URL+"/api/v1/group-test", testData)
	checkResultNoError(t, svr.Client(), svr.URL+"/no-group-test", testData)
	checkResultError(t, svr.Client(), svr.URL+"/error")

	manuallyRegisteredRoutesCount := int32(3)
	assert.GreaterOrEqual(t, r.routesCount.Load(), manuallyRegisteredRoutesCount)
}

func getTestHandlerNoError(testData string) handlers.HTTPHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testData))
		return err
	}
}

func getTestHandlerError(errorData string) handlers.HTTPHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return errors.New(errorData)
	}
}

func checkResultNoError(t *testing.T, client *http.Client, URL string, expectedResult string) {
	resp, err := client.Get(URL)
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if assert.NoError(t, err, "should return response") {
		responseData, err2 := io.ReadAll(resp.Body)
		if assert.NoError(t, err2, "body reading error") {
			assert.Equal(t, expectedResult, string(responseData))
		}
	}
}

func checkResultError(t *testing.T, client *http.Client, URL string) {
	resp, err := client.Get(URL)
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if assert.NoError(t, err, "should return response") {
		responseData, err2 := io.ReadAll(resp.Body)
		if assert.NoError(t, err2, "body reading error") {
			assert.JSONEq(t, `{"errors":[{"message":"Internal Server Error"}]}`, string(responseData))
		}
	}
}
