package cover

import (
	"context"
	"io"
	"log/slog"
)

type BlobStore interface {
	CoverExists(ctx context.Context, filePath string) bool
	GetBookCover(ctx context.Context, filePath string) (io.Reader, error)
}

type Service struct {
	logger    *slog.Logger
	blobStore BlobStore
}

func NewService(logger *slog.Logger, store BlobStore) *Service {
	return &Service{
		logger:    logger,
		blobStore: store,
	}
}

func (s *Service) GetBookCover(ctx context.Context, filePath string) (io.Reader, error) {
	if !s.blobStore.CoverExists(ctx, filePath) {
		return nil, ErrNotFound
	}

	return s.blobStore.GetBookCover(ctx, filePath)
}
