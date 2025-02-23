package handlers

import (
	"context"
	"net/http"
)

type HTTPHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type Middleware func(handler HTTPHandler) HTTPHandler

type RouteRegistrar interface {
	RegisterRoute(method string, group string, path string, handler HTTPHandler, mw ...Middleware)
}
