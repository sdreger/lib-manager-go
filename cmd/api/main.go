package main

import (
	"context"
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/sdreger/lib-manager-go/internal/database"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
)

func main() {
	// ==================== Initialize Logging ====================
	minLogLevel := slog.LevelDebug
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: minLogLevel}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error("Fatal application error", "error", err.Error(), "trace", trace)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) (err error) {
	logger.Info("init service", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// ==================== Configuration Parsing ====================
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
	defer func() {
		logger.Info("API service shutdown complete")
	}()

	// ==================== Get Main Context ====================
	mainCtx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelFunc()

	// ==================== Open DB Connection ====================
	db, err := database.Open(appConfig.DB)
	if err != nil {
		return err
	}
	logger.Info("database connection established", "host", appConfig.DB.Host, "stats", db.Stats())
	defer func() {
		logger.Info("closing database connection")
		if err := db.Close(); err != nil {
			logger.Error("can not close database connection", "error", err)
		}
		logger.Info("database connection closed successfully")
	}()

	// ==================== Run DB Migration ====================
	if err := database.Migrate(logger, appConfig.DB, db.DB); err != nil {
		return err
	}

	// ==================== Start HTTP Server ====================
	if serverAppErr := NewServerApp(appConfig, logger, db).Serve(mainCtx); serverAppErr != nil {
		return serverAppErr
	}

	return nil
}
