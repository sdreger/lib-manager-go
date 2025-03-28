package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers/system"
	handlersV1 "github.com/sdreger/lib-manager-go/cmd/api/handlers/v1"
	"github.com/sdreger/lib-manager-go/internal/blobtstore"
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/sdreger/lib-manager-go/internal/database"
	"github.com/sdreger/lib-manager-go/internal/middleware"
	"log/slog"
	"net/http"
	"sync/atomic"
)

type Router struct {
	mux         *http.ServeMux
	logger      *slog.Logger
	httpConfig  config.HTTPConfig
	routesCount atomic.Int32
	mw          []handlers.Middleware
}

func NewRouter(logger *slog.Logger, db *sqlx.DB, blobStore *blobtstore.MinioStore,
	httpConfig config.HTTPConfig) *Router {

	router := Router{
		mux:         http.NewServeMux(),
		logger:      logger,
		httpConfig:  httpConfig,
		routesCount: atomic.Int32{},
		mw:          []handlers.Middleware{},
	}

	router.registerApplicationMiddlewares()
	router.registerRouteHandlers(db, blobStore)
	logger.Info("router initialized", "registeredRoutes", router.routesCount.Load())

	return &router
}

func (router *Router) GetHandler() http.Handler {
	return router.mux
}

// registerApplicationMiddlewares - register application-wide middlewares.
// Those will be executed first for all registered endpoints, before handler-specific middlewares, and handler itself
// [appMiddleware] -> ... -> [appMiddleware] -> [handlerMiddleware] -> ... -> [handlerMiddleware] -> [handler]
func (router *Router) registerApplicationMiddlewares() {
	// the order matters, first registered - first executed
	router.AddApplicationMiddleware(middleware.Cors(router.httpConfig))
	router.AddApplicationMiddleware(middleware.Errors(router.logger))
	router.AddApplicationMiddleware(middleware.Panics())
}

// registerRouteHandlers - init REST controllers, and delegate route handlers registration to them
func (router *Router) registerRouteHandlers(db *sqlx.DB, blobStore *blobtstore.MinioStore) {
	logger := router.logger
	// the custom DB data type is only needed for system controller to perform health checks
	system.NewController(logger, (*database.DB)(db), blobStore).RegisterRoutes(router)
	handlersV1.NewBookController(logger, db).RegisterRoutes(router)
	handlersV1.NewCoverController(logger, blobStore).RegisterRoutes(router)
	handlersV1.NewFileTypeController(logger, db).RegisterRoutes(router)
}

func (router *Router) AddApplicationMiddleware(mw handlers.Middleware) {
	router.mw = append(router.mw, mw)
}

// RegisterRoute - registers a route handler (endpoint). The group is optional (e.g. '/v1', '/v2', etc.)
// Handler-specific middlewares order matters, first passed - first executed
func (router *Router) RegisterRoute(method string, group string, path string, handler handlers.HTTPHandler,
	mw ...handlers.Middleware) {

	// wrap handler specific middleware around handler being registered
	handler = wrapMiddlewares(handler, mw...)

	// add the application-wide middlewares
	handler = wrapMiddlewares(handler, router.mw...)

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context() // to be able to inject values

		if err := handler(ctx, w, r); err != nil {
			// just a safeguard, because the error should already be handled by the middleware
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			router.logger.Error(err.Error(), slog.Group("request", "method", r.Method, "url", r.URL.String()))
		}
	}

	pattern := fmt.Sprintf("%s %s%s", method, group, path)

	router.routesCount.Add(1)
	router.mux.HandleFunc(pattern, h)
}

func wrapMiddlewares(handler handlers.HTTPHandler, middlewares ...handlers.Middleware) handlers.HTTPHandler {
	// first middleware of the slice is the first to be executed
	for i := len(middlewares) - 1; i >= 0; i-- {
		mw := middlewares[i]
		if mw != nil {
			handler = mw(handler)
		}
	}

	return handler
}
