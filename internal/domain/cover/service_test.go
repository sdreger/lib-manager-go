package cover

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"os"
	"testing"
)

func TestMockBlobStore_GetBoolCover(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))
	ctx := context.Background()
	filePathExists := "publisher/exists.svg"
	existingContent := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>`

	mockBlobStore := NewMockBlobStore(t)
	mockBlobStore.EXPECT().CoverExists(ctx, filePathExists).Return(true).Once()
	mockBlobStore.EXPECT().GetBookCover(ctx, filePathExists).
		Return(bytes.NewBufferString(existingContent), nil).Once()
	service := NewService(logger, mockBlobStore)

	cover, err := service.GetBookCover(ctx, filePathExists)
	require.NoError(t, err, "should return book cover")
	content, err := io.ReadAll(cover)
	require.NoError(t, err)
	require.Equal(t, existingContent, string(content))
}

func TestMockBlobStore_GetBoolCover_NotExist(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))
	ctx := context.Background()
	filePathNotExist := "publisher/not-exist.svg"

	mockBlobStore := NewMockBlobStore(t)
	mockBlobStore.EXPECT().CoverExists(ctx, filePathNotExist).Return(false).Once()
	service := NewService(logger, mockBlobStore)

	nonExistingCover, err := service.GetBookCover(ctx, filePathNotExist)
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, nonExistingCover)
}
