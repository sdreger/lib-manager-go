package main

import (
	"context"
	"fmt"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"testing"
	"time"
)

var testResponse = `{"status":"OK"}`

func TestServerApp_ServeAndShutdown(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	randomServerPort := getRandomPort()

	appConfig, err := config.New()
	if assert.NoError(t, err, "should create a new default config") {
		appConfig.HTTP.Port = randomServerPort
		serverApp := NewServerApp(appConfig, logger)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		group, gCtx := errgroup.WithContext(ctx)
		group.Go(func() error {
			serverApp.router.RegisterRoute(http.MethodGet, "", "/test", getTestHandler())
			return serverApp.Serve(gCtx)
		})

		time.Sleep(1 * time.Second)
		checkTestHandlerResponse(t, randomServerPort)
		cancel()

		err := group.Wait()
		assert.NoError(t, err, "should shutdown gracefully")
	}
}

func getTestHandler() handlers.HTTPHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(testResponse))
		return err
	}
}

func checkTestHandlerResponse(t *testing.T, serverPort int) {

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/test", serverPort))
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if assert.NoError(t, err, "should return response") {
		responseData, err2 := io.ReadAll(resp.Body)
		if assert.NoError(t, err2, "body reading error") {
			assert.Equal(t, testResponse, string(responseData))
		}
	}
}

func getRandomPort() int {
	maxPort := 65534
	minPort := 1024
	return rand.IntN(maxPort-minPort) + minPort
}
