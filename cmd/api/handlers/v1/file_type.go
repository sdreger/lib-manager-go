package v1

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/domain/filetype"
	"github.com/sdreger/lib-manager-go/internal/paging"
	"github.com/sdreger/lib-manager-go/internal/response"
	"log/slog"
	"net/http"
)

type FileTypeService interface {
	GetFileTypes(
		ctx context.Context,
		pageRequest paging.PageRequest,
		sort paging.Sort,
	) (paging.Page[filetype.LookupItem], error)
}

type FileTypeController struct {
	logger  *slog.Logger
	service FileTypeService
}

func NewFileTypeController(logger *slog.Logger, db *sqlx.DB) *FileTypeController {
	return &FileTypeController{logger: logger, service: filetype.NewService(logger, db)}
}

func (cnt *FileTypeController) RegisterRoutes(registrar handlers.RouteRegistrar) {
	registrar.RegisterRoute(http.MethodGet, group, "/file_types", cnt.GetFileTypes)
}

func (cnt *FileTypeController) GetFileTypes(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, pageErr := paging.NewPageRequest(r.URL.Query())
	if pageErr != nil {
		return pageErr
	}

	sort, sortErr := paging.NewSort(r.URL.Query(), filetype.AllowedSortFields)
	if sortErr != nil {
		return sortErr
	}

	FileTypePage, err := cnt.service.GetFileTypes(ctx, page, sort)
	if err != nil {
		return err
	}

	return response.RenderDataJSON(w, http.StatusOK, FileTypePage)
}
