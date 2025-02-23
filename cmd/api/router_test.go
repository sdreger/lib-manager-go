package main

import (
	"context"
	"errors"
	apiErrors "github.com/sdreger/lib-manager-go/cmd/api/errors"
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

	r := NewRouter(logger, nil)
	r.RegisterRoute(http.MethodGet, "/api/v1", "/group-test", getTestHandlerNoError(testData))
	r.RegisterRoute(http.MethodGet, "", "/no-group-test", getTestHandlerNoError(testData))
	r.RegisterRoute(http.MethodGet, "", "/error", getTestHandlerError())
	r.RegisterRoute(http.MethodGet, "", "/not-found-error", getTestHandlerNotFoundError())
	r.RegisterRoute(http.MethodGet, "", "/validation-error", getTestHandlerValidationError())
	r.RegisterRoute(http.MethodGet, "", "/validation-errors", getTestHandlerValidationErrors())

	svr := httptest.NewServer(r.GetHandler())
	defer svr.Close()
	checkResultNoError(t, svr.Client(), svr.URL+"/api/v1/group-test", testData)
	checkResultNoError(t, svr.Client(), svr.URL+"/no-group-test", testData)
	checkResultError(t, svr.Client(), svr.URL+"/error")
	checkResultNotFoundError(t, svr.Client(), svr.URL+"/not-found-error")
	checkResultValidationError(t, svr.Client(), svr.URL+"/validation-error")
	checkResultValidationErrors(t, svr.Client(), svr.URL+"/validation-errors")

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

func getTestHandlerError() handlers.HTTPHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return errors.New("internal error")
	}
}

func getTestHandlerNotFoundError() handlers.HTTPHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return apiErrors.ErrNotFound
	}
}

func getTestHandlerValidationError() handlers.HTTPHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return apiErrors.ValidationError{
			Field:   "title",
			Message: "title is required",
		}
	}
}

func getTestHandlerValidationErrors() handlers.HTTPHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return apiErrors.ValidationErrors{
			apiErrors.ValidationError{
				Field:   "name",
				Message: "name is required",
			},
			apiErrors.ValidationError{
				Field:   "publisher",
				Message: "publisher can not be empty",
			},
		}
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
			assert.JSONEq(t, `{"errors":[{"message":"Internal Server Error"}]}`, string(response))
		}
	}
}

func checkResultNotFoundError(t *testing.T, client *http.Client, URL string) {
	resp, err := client.Get(URL)
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if assert.NoError(t, err, "should return response") {
		response, err2 := io.ReadAll(resp.Body)
		if assert.NoError(t, err2, "body reading error") {
			assert.JSONEq(t, `{"errors":[{"message":"the requested resource could not be found"}]}`, string(response))
		}
	}
}

func checkResultValidationError(t *testing.T, client *http.Client, URL string) {
	resp, err := client.Get(URL)
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if assert.NoError(t, err, "should return response") {
		response, err2 := io.ReadAll(resp.Body)
		if assert.NoError(t, err2, "body reading error") {
			assert.JSONEq(t, `{"errors":[{"field":"title", "message":"title is required"}]}`, string(response))
		}
	}
}

func checkResultValidationErrors(t *testing.T, client *http.Client, URL string) {
	resp, err := client.Get(URL)
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if assert.NoError(t, err, "should return response") {
		response, err2 := io.ReadAll(resp.Body)
		if assert.NoError(t, err2, "body reading error") {
			assert.JSONEq(t,
				`{"errors":[{"field":"name", "message":"name is required"},
						  			 {"field":"publisher", "message":"publisher can not be empty"}]}`,
				string(response))
		}
	}
}
