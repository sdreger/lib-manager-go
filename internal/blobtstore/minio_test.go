package blobtstore

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	testMinio "github.com/testcontainers/testcontainers-go/modules/minio"
	"io"
	"log"
	"log/slog"
	"os"
	"testing"
)

const (
	testSVG = `
		<?xml version="1.0" encoding="UTF-8" standalone="no"?>
		<svg xmlns="http://www.w3.org/2000/svg" width="500" height="500">
		<circle cx="250" cy="250" r="210" fill="#fff" stroke="#000" stroke-width="8"/>
		</svg>
		`
)

func TestNewMinioStore(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))

	minioContainer, err := testMinio.Run(ctx, "quay.io/minio/minio:RELEASE.2025-02-18T16-25-55Z")
	defer func() {
		if err := testcontainers.TerminateContainer(minioContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	require.NoError(t, err, "failed to start minio container")

	endpoint, err := minioContainer.ConnectionString(ctx)
	require.NoError(t, err, "failed to get container endpoint")

	appConfig, err := config.New()
	require.NoError(t, err, "failed to load app config")
	blobStoreConfig := appConfig.BLOBStore
	blobStoreConfig.MinioEndpoint = endpoint
	blobStoreConfig.MinioAccessKeyID = minioContainer.Username
	blobStoreConfig.MinioSecretAccessKey = minioContainer.Password

	minioStore, err := NewMinioStore(logger, blobStoreConfig)
	require.NoError(t, err, "failed to create minio store")

	err = minioStore.CreateBuckets(ctx)
	require.NoError(t, err, "failed to create buckets")

	t.Run("GetBookCover", func(t *testing.T) {
		coverPath := "publisher/test_file.svg"
		err := storeBookCover(ctx, minioStore, coverPath, testSVG)
		require.NoError(t, err, "failed to store book cover")

		require.True(t, minioStore.CoverExists(ctx, coverPath))

		bookCoverStub, err := minioStore.GetBookCover(ctx, coverPath)
		require.NoError(t, err, "failed to get book cover")

		fileContent, err := io.ReadAll(bookCoverStub)
		require.NoError(t, err, "failed to read book cover")

		require.Equal(t, []byte(testSVG), fileContent)
	})

	t.Run("CoverDoesNotExist", func(t *testing.T) {
		coverPath := "publisher/wrong_file.svg"
		assert.False(t, minioStore.CoverExists(ctx, coverPath))
	})

	t.Run("BucketAlreadyExists", func(t *testing.T) {
		err = minioStore.CreateBuckets(ctx)
		require.NoError(t, err, "bucket creation is idempotent")
	})

	t.Run("WrongBucketName", func(t *testing.T) {
		minioStore.coverBucketName = "a"
		err = minioStore.CreateBuckets(ctx)
		require.ErrorContains(t, err, "Bucket name cannot be shorter than 3 characters", "bucket name is invalid")
	})
}

func TestNewMinioStore_CanNotCreateClient(t *testing.T) {
	appConfig, err := config.New()
	require.NoError(t, err, "failed to load app config")
	blobStoreConfig := appConfig.BLOBStore
	blobStoreConfig.MinioEndpoint = ""

	_, err = NewMinioStore(nil, blobStoreConfig)
	require.ErrorContains(t, err, "does not follow ip address or domain name standards")
}

func storeBookCover(ctx context.Context, minioStore *MinioStore, path string, content string) error {
	contentBuffer := bytes.NewBufferString(content)
	_, err := minioStore.client.PutObject(ctx, minioStore.config.BookCoverBucket,
		path, contentBuffer, int64(contentBuffer.Len()), minio.PutObjectOptions{})

	return err
}
