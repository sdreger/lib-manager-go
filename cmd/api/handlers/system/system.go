package system

import (
	"context"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/response"
	"log/slog"
	"net/http"
)

const (
	statusFail healthCheckStatus = "FAIL"
	statusOK   healthCheckStatus = "OK"
)

type HealthChecker interface {
	HealthCheck(ctx context.Context) error
	HealthCheckID() string
}

type healthCheckStatus string

type healthCheckResult struct {
	id    string
	error error
}

type Controller struct {
	logger         *slog.Logger
	healthCheckers []HealthChecker
}

func NewController(logger *slog.Logger, healthCheckers ...HealthChecker) *Controller {
	return &Controller{
		logger:         logger,
		healthCheckers: healthCheckers,
	}
}

func (cnt *Controller) RegisterRoutes(registrar handlers.RouteRegistrar) {
	registrar.RegisterRoute(http.MethodGet, "", "/livez", cnt.LivenessProbe)
	registrar.RegisterRoute(http.MethodGet, "", "/readyz", cnt.ReadinessProbe)
}

func (cnt *Controller) LivenessProbe(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	defer cnt.closeBody(req)

	data := map[string]string{
		"status": "OK",
	}

	return response.RenderJSONWithHeaders(w, http.StatusOK, data, nil)
}

func (cnt *Controller) ReadinessProbe(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	defer cnt.closeBody(req)

	resultChan := make(chan healthCheckResult, len(cnt.healthCheckers))
	data := make(map[string]healthCheckStatus, len(cnt.healthCheckers))
	for _, healthChecker := range cnt.healthCheckers {
		data[healthChecker.HealthCheckID()] = statusFail
		go func(hc HealthChecker) {
			resultChan <- healthCheckResult{
				id:    hc.HealthCheckID(),
				error: hc.HealthCheck(ctx),
			}
		}(healthChecker)
	}

	responseStatus := http.StatusOK
	for i := 0; i < len(cnt.healthCheckers); i++ {
		result := <-resultChan
		if result.error != nil {
			responseStatus = http.StatusInternalServerError
			cnt.logger.Error("health check failed: "+result.error.Error(), "healthCheckID", result.id)
		} else {
			data[result.id] = statusOK
		}
	}
	close(resultChan)

	return response.RenderJSONWithHeaders(w, responseStatus, data, nil)
}

func (cnt *Controller) closeBody(req *http.Request) {
	if req.Body != nil {
		_ = req.Body.Close()
	}
}
