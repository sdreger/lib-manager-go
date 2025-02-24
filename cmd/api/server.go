package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sdreger/lib-manager-go/internal/config"
	"log/slog"
	"net/http"
)

type ServerApp struct {
	config config.AppConfig
	logger *slog.Logger
	router *Router
}

func NewServerApp(config config.AppConfig, logger *slog.Logger, db *sqlx.DB) *ServerApp {
	return &ServerApp{
		config: config,
		logger: logger,
		router: NewRouter(logger, db, config.HTTP),
	}
}

func (app *ServerApp) Serve(ctx context.Context) error {
	appConfig := app.config
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", appConfig.HTTP.Host, appConfig.HTTP.Port),
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelWarn),
		ReadTimeout:  appConfig.HTTP.ReadTimeout,
		WriteTimeout: appConfig.HTTP.WriteTimeout,
		IdleTimeout:  appConfig.HTTP.IdleTimeout,
		Handler:      app.router.GetHandler(),
	}
	shutdownErrorChan := make(chan error, 1)
	go func() {
		// Wait for notifyContext is closed (for the graceful shutdown)
		<-ctx.Done()
		app.logger.Info("graceful server shutdown initiated")
		shutdownCtx, cancelFunc := context.WithTimeout(context.Background(), appConfig.HTTP.ShutdownTimeout)
		defer cancelFunc()
		shutdownErrorChan <- server.Shutdown(shutdownCtx)
	}()

	app.logger.Info("starting server",
		slog.Group("server", "host", appConfig.HTTP.Host, "port", appConfig.HTTP.Port))
	serverError := server.ListenAndServe()
	if !errors.Is(serverError, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", serverError)
	}

	// ==================== Server Shutdown ====================
	err := <-shutdownErrorChan
	if err != nil {
		return fmt.Errorf("graceful server shutdown error: %w", err)
	}
	app.logger.Info("graceful server shutdown complete")

	return nil
}
