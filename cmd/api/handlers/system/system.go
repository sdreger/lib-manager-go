package system

import (
	"context"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/response"
	"log/slog"
	"net/http"
)

type Handler struct {
	logger *slog.Logger
}

func RegisterHandler(logger *slog.Logger, registrar handlers.RouteRegistrar) {
	handler := Handler{
		logger: logger,
	}

	registrar.RegisterRoute(http.MethodGet, "", "/health", handler.HealthProbe)
}

func (h Handler) HealthProbe(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	defer func() {
		if req.Body != nil {
			req.Body.Close()
		}
	}()

	data := map[string]string{
		"status": "OK",
	}

	return response.RenderJSONWithHeaders(w, http.StatusOK, data, nil)
}
