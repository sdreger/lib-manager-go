package middleware

import (
	"context"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/config"
	"net/http"
	"strings"
)

const (
	headerOrigin         = "Origin"
	headerAllowedOrigin  = "Access-Control-Allow-Origin"
	headerAllowedMethods = "Access-Control-Allow-Methods"
	headerAllowedHeaders = "Access-Control-Allow-Headers"
)

func Cors(config config.HTTPConfig) handlers.Middleware {
	return func(next handlers.HTTPHandler) handlers.HTTPHandler {
		if len(config.CORS.AllowedOrigins) == 0 {
			return next
		}

		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			origin := r.Header.Get(headerOrigin)
			for _, allowedOrigin := range config.CORS.AllowedOrigins {
				if strings.Contains(origin, allowedOrigin) {
					w.Header().Set(headerAllowedOrigin, origin)
					w.Header().Set(headerAllowedMethods, strings.Join(config.CORS.AllowedMethods, ","))
					w.Header().Set(headerAllowedHeaders, strings.Join(config.CORS.AllowedHeaders, ","))
				}
			}

			return next(ctx, w, r)
		}
	}
}
