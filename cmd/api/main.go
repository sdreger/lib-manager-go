package main

import (
	"context"
	"github.com/sdreger/lib-manager-go/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
)

func main() {
	// ==================== Logging ====================
	minLogLevel := slog.LevelDebug
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: minLogLevel}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) (err error) {
	logger.Info("init service", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// ==================== Configuration parsing ====================
	appConfig, err := config.New()
	if err != nil {
		return err
	}

	// ==================== Start API Service ====================
	logger.Info("starting API service", slog.Group("build",
		"revision", appConfig.BuildInfo.Revision,
		"time", appConfig.BuildInfo.Time,
		"dirty", appConfig.BuildInfo.Dirty,
	))

	mainCtx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelFunc()

	serverAppErr := NewServerApp(appConfig, logger).Serve(mainCtx)
	if serverAppErr != nil {
		logger.Error("API service error", "error", serverAppErr.Error())
	}

	logger.Info("API service shutdown complete")
	return nil
}
