package middleware

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPanics_NoPanic(t *testing.T) {
	expectedError := errors.New("some error")
	middleware := Panics()
	handler := middleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return expectedError
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	err := handler(context.Background(), recorder, request)
	require.Error(t, err)
	require.ErrorIs(t, err, expectedError)
}

func TestPanics_Panic(t *testing.T) {
	expectedError := errors.New("some error")
	middleware := Panics()
	handler := middleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		panic(expectedError)
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	err := handler(context.Background(), recorder, request)
	require.Error(t, err)
	require.ErrorContains(t, err, "recover from panic: some error")
}
