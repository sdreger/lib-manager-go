package system

import (
	"context"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/response"
	"log/slog"
	"net/http"
)

type Controller struct {
	logger *slog.Logger
}

func NewController(logger *slog.Logger) *Controller {
	return &Controller{logger: logger}
}

func (cnt *Controller) RegisterRoutes(registrar handlers.RouteRegistrar) {
	registrar.RegisterRoute(http.MethodGet, "", "/health", cnt.HealthProbe)
}

func (cnt *Controller) HealthProbe(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
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
