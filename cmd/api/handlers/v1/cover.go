package v1

import (
	"context"
	"errors"
	"fmt"
	apiErrors "github.com/sdreger/lib-manager-go/cmd/api/errors"
	"github.com/sdreger/lib-manager-go/cmd/api/handlers"
	"github.com/sdreger/lib-manager-go/internal/blobtstore"
	"github.com/sdreger/lib-manager-go/internal/domain/cover"
	"io"
	"log/slog"
	"net/http"
)

const (
	publisherNamePathVariable = "publisherName"
	coverFileNamePathVariable = "coverFileName"
)

type CoverService interface {
	GetBookCover(ctx context.Context, filePath string) (io.Reader, error)
}

type CoverHandler struct {
	logger       *slog.Logger
	coverService CoverService
}

func NewCoverHandler(logger *slog.Logger, blobStore *blobtstore.MinioStore) *CoverHandler {
	return &CoverHandler{
		logger:       logger,
		coverService: cover.NewService(logger, blobStore),
	}
}

func (ch *CoverHandler) RegisterHandler(registrar handlers.RouteRegistrar) {
	group := "/v1"
	registrar.RegisterRoute(http.MethodGet, group, "/covers/{publisherName}/{coverFileName}", ch.GetBookCover)
}

func (ch *CoverHandler) GetBookCover(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	publisherName := r.PathValue(publisherNamePathVariable)
	coverFileName := r.PathValue(coverFileNamePathVariable)
	bookCoverReader, err := ch.coverService.GetBookCover(ctx, fmt.Sprintf("%s/%s", publisherName, coverFileName))
	if errors.Is(err, cover.ErrNotFound) {
		return apiErrors.ErrNotFound
	}
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "octet/stream")
	_, coverWriteError := io.Copy(w, bookCoverReader)

	return coverWriteError
}
