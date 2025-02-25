package middleware

import (
	"context"
	"errors"
	apiErrors "github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/response"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Errors(logger *slog.Logger) handlers.Middleware {
	return func(next handlers.HTTPHandler) handlers.HTTPHandler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := next(ctx, w, r); err != nil {
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

				reportServerError(logger, r, err, unexpectedError)
			}

			// the error has been handled, no need to propagate it
			return nil
		}
	}
}

func reportServerError(logger *slog.Logger, r *http.Request, err error, isUnexpected bool) {
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
	logger.Error(message, requestGroup)
}
