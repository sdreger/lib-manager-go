package spec

import (
	"context"
	"embed"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"gopkg.org/swaggerui"
	"log/slog"
	"net/http"
)

const (
	specURLPrefix = "/v3/api-docs/"
	specFileName  = "openapi.yaml"
)

//go:embed openapi.yaml
var swaggerUI embed.FS

type Controller struct {
	logger *slog.Logger
}

func NewController(logger *slog.Logger) *Controller {
	return &Controller{
		logger: logger,
	}
}

func (cnt *Controller) RegisterRoutes(registrar handlers.RouteRegistrar) {
	registrar.RegisterRoute(http.MethodGet, "", "/swagger-ui/", cnt.SwaggerUI)
	registrar.RegisterRoute(http.MethodGet, "", "/v3/api-docs/", cnt.OpenAPISpec)
}

func (cnt *Controller) SwaggerUI(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	defer cnt.closeBody(req)

	swaggerURL := specURLPrefix + specFileName
	http.StripPrefix("/swagger-ui/", swaggerui.Handler(swaggerURL)).ServeHTTP(w, req)

	return nil
}

func (cnt *Controller) OpenAPISpec(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	defer cnt.closeBody(req)

	http.StripPrefix(specURLPrefix, http.FileServer(http.FS(swaggerUI))).ServeHTTP(w, req)

	return nil
}

func (cnt *Controller) closeBody(req *http.Request) {
	if req.Body != nil {
		_ = req.Body.Close()
	}
}
