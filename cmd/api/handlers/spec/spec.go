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
	specURLPrefix    = "/v3/api-docs/"
	specFileName     = "openapi.yaml"
	swaggerURLPrefix = "/swagger-ui/"
)

//go:embed openapi.yaml
var apiSpec embed.FS

type Controller struct {
	logger *slog.Logger
}

func NewController(logger *slog.Logger) *Controller {
	return &Controller{
		logger: logger,
	}
}

func (cnt *Controller) RegisterRoutes(registrar handlers.RouteRegistrar) {
	registrar.RegisterRoute(http.MethodGet, "", swaggerURLPrefix, cnt.SwaggerUI)
	registrar.RegisterRoute(http.MethodGet, "", specURLPrefix, cnt.APISpec)
}

func (cnt *Controller) SwaggerUI(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	defer cnt.closeBody(req)

	specURL := specURLPrefix + specFileName
	http.StripPrefix(swaggerURLPrefix, swaggerui.Handler(specURL)).ServeHTTP(w, req)

	return nil
}

func (cnt *Controller) APISpec(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	defer cnt.closeBody(req)

	http.StripPrefix(specURLPrefix, http.FileServer(http.FS(apiSpec))).ServeHTTP(w, req)

	return nil
}

func (cnt *Controller) closeBody(req *http.Request) {
	if req.Body != nil {
		_ = req.Body.Close()
	}
}
