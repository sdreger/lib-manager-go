package middleware

import (
	"context"
	"fmt"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"net/http"
)

// Panics - middleware for capturing panics and converting them to errors.
// Should be registered as the last application-wide middleware, before any handler-specific middlewares
func Panics() handlers.Middleware {
	return func(next handlers.HTTPHandler) handlers.HTTPHandler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			// recover from a potential panic and set the err value
			defer func() {
				if rec := recover(); rec != nil {
					err = fmt.Errorf("recover from panic: %v", rec)
				}
			}()

			return next(ctx, w, r)
		}
	}
}
