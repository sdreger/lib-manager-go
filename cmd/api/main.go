package main

import (
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
	logger.Info("Application init", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	
	return nil
}
