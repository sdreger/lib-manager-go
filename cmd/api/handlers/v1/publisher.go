package v1

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/domain/publisher"
	"github.com/sdreger/lib-manager-go/internal/paging"
	"github.com/sdreger/lib-manager-go/internal/response"
	"log/slog"
	"net/http"
)

type PublisherService interface {
	GetPublishers(
		ctx context.Context,
		pageRequest paging.PageRequest,
		sort paging.Sort,
	) (paging.Page[publisher.LookupItem], error)
}

type PublisherController struct {
	logger  *slog.Logger
	service PublisherService
}

func NewPublisherController(logger *slog.Logger, db *sqlx.DB) *PublisherController {
	return &PublisherController{logger: logger, service: publisher.NewService(logger, db)}
}

func (cnt *PublisherController) RegisterRoutes(registrar handlers.RouteRegistrar) {
	registrar.RegisterRoute(http.MethodGet, group, "/publishers", cnt.GetPublishers)
}

func (cnt *PublisherController) GetPublishers(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, pageErr := paging.NewPageRequest(r.URL.Query())
	if pageErr != nil {
		return pageErr
	}

	sort, sortErr := paging.NewSort(r.URL.Query(), publisher.AllowedSortFields)
	if sortErr != nil {
		return sortErr
	}

	publisherPage, err := cnt.service.GetPublishers(ctx, page, sort)
	if err != nil {
		return err
	}

	return response.RenderDataJSON(w, http.StatusOK, publisherPage)
}
