package middleware

import (
	"context"
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCors_NoHeaders(t *testing.T) {
	appConfig, err := config.New()
	require.NoError(t, err)

	middleware := Cors(appConfig.HTTP)
	handler := middleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return nil
	})
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	err = handler(context.Background(), response, request)
	require.NoError(t, err)
	require.Equal(t, "", response.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, "", response.Header().Get("Access-Control-Allow-Headers"))
	require.Equal(t, "", response.Header().Get("Access-Control-Allow-Methods"))
}

func TestCors_Headers(t *testing.T) {
	appConfig, err := config.New()
	require.NoError(t, err)
	allowedOrigin01 := "http://127.0.0.1:3000"
	allowedOrigin02 := "http://192.168.0.10:3000"
	allowedMethod01 := "GET"
	allowedMethod02 := "POST"
	allowedHeader01 := "X-Test-Header"
	allowedHeader02 := "Content-Type"
	appConfig.HTTP.CORS.AllowedOrigins = []string{allowedOrigin01, allowedOrigin02}
	appConfig.HTTP.CORS.AllowedMethods = []string{allowedMethod01, allowedMethod02}
	appConfig.HTTP.CORS.AllowedHeaders = []string{allowedHeader01, allowedHeader02}

	middleware := Cors(appConfig.HTTP)
	handler := middleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return nil
	})
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Add("Origin", allowedOrigin01)
	response := httptest.NewRecorder()

	err = handler(context.Background(), response, request)
	require.NoError(t, err)
	require.Equal(t, allowedOrigin01, response.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, strings.Join(appConfig.HTTP.CORS.AllowedHeaders, ","),
		response.Header().Get("Access-Control-Allow-Headers"))
	require.Equal(t, strings.Join(appConfig.HTTP.CORS.AllowedMethods, ","),
		response.Header().Get("Access-Control-Allow-Methods"))
}
