package main

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	apiErrors "github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers/system"
	handlersV1 "github.com/sdreger/lib-manager-go/cmd/api/handlers/v1"
	"github.com/sdreger/lib-manager-go/internal/response"
	"log/slog"
	"net/http"
	"runtime/debug"
	"sync/atomic"
)

type Router struct {
	mux         *http.ServeMux
	logger      *slog.Logger
	routesCount atomic.Int32
	mw          []handlers.Middleware
}

func NewRouter(logger *slog.Logger, db *sqlx.DB) *Router {
	router := new(Router).
		WithMux(http.NewServeMux()).
		WithLogger(logger)

	router.registerHandlers(db)
	logger.Info("router initialized", "registeredRoutes", router.routesCount.Load())

	return router
}

func (router *Router) WithMux(mux *http.ServeMux) *Router {
	router.mux = mux
	return router
}

func (router *Router) WithLogger(logger *slog.Logger) *Router {
	router.logger = logger
	return router
}

func (router *Router) WithMiddleware(mw handlers.Middleware) *Router {
	router.mw = append(router.mw, mw)
	return router
}

func (router *Router) GetHandler() http.Handler {
	return router.mux
}

func (router *Router) registerHandlers(db *sqlx.DB) {
	system.NewHandler(router.logger).RegisterHandler(router)
	handlersV1.NewBookHandler(router.logger, db).RegisterHandler(router)
}

func (router *Router) RegisterRoute(method string, group string, path string, handler handlers.HTTPHandler,
	mw ...handlers.Middleware) {

	// wrap handler specific middleware around handler being registered
	handler = wrapMiddlewares(handler, mw...)

	// add the application-wide middlewares
	handler = wrapMiddlewares(handler, router.mw...)

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context() // to be able to inject values
		if err := handler(ctx, w, r); err != nil {
			router.handleServerError(w, r, err)
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

func (router *Router) handleServerError(w http.ResponseWriter, r *http.Request, err error) {
	var validationError apiErrors.ValidationError
	var validationErrors apiErrors.ValidationErrors
	var renderingError error
	var unexpectedError bool
	switch {
	case errors.As(err, &validationError):
		renderingError = response.RenderErrorJSON(w, http.StatusBadRequest,
			[]response.APIError{validationError.ToAPIError()})
	case errors.As(err, &validationErrors):
		renderingError = response.RenderErrorJSON(w, http.StatusBadRequest,
			validationErrors.ToAPIErrors())
	case errors.Is(err, apiErrors.ErrNotFound):
		renderingError = response.RenderErrorJSON(w, http.StatusNotFound,
			[]response.APIError{{Message: err.Error()}})
	default:
		renderingError = response.RenderErrorJSON(w, http.StatusInternalServerError,
			[]response.APIError{{Message: http.StatusText(http.StatusInternalServerError)}})
		unexpectedError = true
	}

	if renderingError != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	router.reportServerError(r, err, unexpectedError)
}

func (router *Router) reportServerError(r *http.Request, err error, isUnexpected bool) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
	)

	var requestGroup slog.Attr
	if isUnexpected {
		trace := string(debug.Stack())
		requestGroup = slog.Group("request", "method", method, "url", url, "trace", trace)
	} else {
		requestGroup = slog.Group("request", "method", method, "url", url)
	}
	router.logger.Error(message, requestGroup)
}
