package middleware

import (
	"context"
	"errors"
	apiErrors "github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestErrors_ValidationError(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	middleware := Errors(logger)
	handler := middleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return apiErrors.ValidationError{
			Field:   "title",
			Message: "title is required",
		}
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	err := handler(context.Background(), recorder, request)
	require.NoError(t, err, "error should be handled by middleware")
	require.Equal(t, http.StatusBadRequest, recorder.Code)

	body, err := io.ReadAll(recorder.Result().Body)
	if assert.NoError(t, err, "body reading error") {
		assert.JSONEq(t, `{"errors":[{"field":"title", "message":"title is required"}]}`, string(body))
	}
}

func TestErrors_ValidationErrors(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	middleware := Errors(logger)
	handler := middleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
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
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	err := handler(context.Background(), recorder, request)
	require.NoError(t, err, "error should be handled by middleware")
	require.Equal(t, http.StatusBadRequest, recorder.Code)

	body, err := io.ReadAll(recorder.Result().Body)
	if assert.NoError(t, err, "body reading error") {
		assert.JSONEq(t, `{"errors":[{"field":"name", "message":"name is required"},
								   {"field":"publisher", "message":"publisher can not be empty"}]}`, string(body))
	}
}

func TestErrors_NotFoundError(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	middleware := Errors(logger)
	handler := middleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return apiErrors.ErrNotFound
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	err := handler(context.Background(), recorder, request)
	require.NoError(t, err, "error should be handled by middleware")
	require.Equal(t, http.StatusNotFound, recorder.Code)

	body, err := io.ReadAll(recorder.Result().Body)
	if assert.NoError(t, err, "body reading error") {
		assert.JSONEq(t, `{"errors":[{"message":"the requested resource could not be found"}]}`, string(body))
	}
}

func TestErrors_UnexpectedError(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	middleware := Errors(logger)
	handler := middleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return errors.New("unexpected internal error")
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	err := handler(context.Background(), recorder, request)
	require.NoError(t, err, "error should be handled by middleware")
	require.Equal(t, http.StatusInternalServerError, recorder.Code)

	body, err := io.ReadAll(recorder.Result().Body)
	if assert.NoError(t, err, "body reading error") {
		assert.JSONEq(t, `{"errors":[{"message":"Internal Server Error"}]}`, string(body))
	}
}

func TestErrors_RenderingError(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	middleware := Errors(logger)
	handler := middleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return errors.New("unexpected internal error")
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := ErroredResponseRecorder{}
	err := handler(context.Background(), &recorder, request)
	require.NoError(t, err, "error should be handled by middleware")
	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}

type ErroredResponseRecorder struct {
	httptest.ResponseRecorder
}

func (rec *ErroredResponseRecorder) Write(buf []byte) (int, error) {
	if string(buf) == `{"errors":[{"message":"Internal Server Error"}]}` {
		return 0, errors.New("write error")
	}

	return rec.ResponseRecorder.Write(buf)
}
