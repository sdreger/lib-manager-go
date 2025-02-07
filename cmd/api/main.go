package main

import (
	"github.com/sdreger/lib-manager-go/internal/config"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
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

func run(logger *slog.Logger) error {
	logger.Info("init service", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	appConfig, err := config.New()
	if err != nil {
		return err
	}

	logger.Info("starting service", slog.Group("build",
		"revision", appConfig.BuildInfo.Revision,
		"time", appConfig.BuildInfo.Time,
		"dirty", appConfig.BuildInfo.Dirty,
	))

	return nil
}
