package main

import (
	"fmt"
	"github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers/system"
	"log/slog"
	"net/http"
	"runtime/debug"
	"sync/atomic"
)

type Router struct {
	mux         *http.ServeMux
	logger      *slog.Logger
	routesCount atomic.Int32
}

func NewRouter(logger *slog.Logger) *Router {
	router := Router{
		mux:    http.NewServeMux(),
		logger: logger,
	}
	router.RegisterHandlers()
	logger.Info("router initialized", "registeredRoutes", router.routesCount.Load())

	return &router
}

func (router *Router) GetHandler() http.Handler {
	return router.mux
}

func (router *Router) RegisterHandlers() {
	system.NewHandler(router.logger).RegisterHandler(router)
}

func (router *Router) RegisterRoute(method string, group string, path string, handler handlers.HTTPHandler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context() // to be able to inject values
		// TODO: add middleware injection here
		if err := handler(ctx, w, r); err != nil {
			router.reportServerError(r, err)
			errors.HandleError(w, err)
		}
	}

	pattern := fmt.Sprintf("%s %s%s", method, group, path)

	router.routesCount.Add(1)
	router.mux.HandleFunc(pattern, h)
}

func (router *Router) reportServerError(r *http.Request, err error) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)

	requestGroup := slog.Group("request", "method", method, "url", url)
	router.logger.Error(message, requestGroup)
	router.logger.Error(message, requestGroup, "trace", trace)
}
