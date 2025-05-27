package blobtstore

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/sdreger/lib-manager-go/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"
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

	minioContainer := tests.StartMinioTestContainer(t)
	minioConfig := tests.GetTestMinioConfig(t, minioContainer)
	minioStore := setUpTestMinioStore(t, minioConfig)

	err := minioStore.CreateBuckets(ctx)
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

func TestMinioStore_HealthCheck(t *testing.T) {
	ctx := context.Background()
	minioContainer := tests.StartMinioTestContainer(t)
	minioConfig := tests.GetTestMinioConfig(t, minioContainer)
	minioStore := setUpTestMinioStore(t, minioConfig)

	err := minioStore.HealthCheck(ctx)
	require.NoError(t, err, "failed to perform a healthcheck")

	tests.TerminateMinioContainer(t, minioContainer)

	// a client call should be done to mark the client as 'offline'
	err = minioStore.CreateBuckets(ctx)
	require.Error(t, err, "should fail to create buckets")

	err = minioStore.HealthCheck(ctx)
	require.Error(t, err, "should fail to perform a healthcheck")

	require.Equal(t, "minio", minioStore.HealthCheckID())
}

func TestMinioStore_HealthCheckStartError(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))
	appConfig, err := config.New()
	require.NoError(t, err, "failed to load app config")

	appConfig.BLOBStore.MinioHealthCheckInterval = time.Duration(100) * time.Millisecond // too small
	_, err = NewMinioStore(logger, appConfig.BLOBStore)
	require.Error(t, err, "failed to create minio store")
}

func setUpTestMinioStore(t *testing.T, blobStoreConfig config.BLOBStoreConfig) *MinioStore {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(100)}))

	minioStore, err := NewMinioStore(logger, blobStoreConfig)
	require.NoError(t, err, "failed to create minio store")
	waitForHealthcheckPass(t, minioStore)

	t.Cleanup(func() {
		minioStore.Close()
	})

	return minioStore
}

func waitForHealthcheckPass(t *testing.T, minioStore *MinioStore) {
	for i := 0; i < 5; i++ {
		if err := minioStore.HealthCheck(t.Context()); err != nil {
			t.Logf("waiting for Minio healthckeck pass...: %v", err)
			time.Sleep(time.Second)
		} else {
			t.Logf("Minio healthcheck passed")
			break
		}
	}

	if err := minioStore.HealthCheck(t.Context()); err != nil {
		t.Fatalf("Minio healthcheck failed: %v", err)
	}
}

func storeBookCover(ctx context.Context, minioStore *MinioStore, path string, content string) error {
	contentBuffer := bytes.NewBufferString(content)
	_, err := minioStore.client.PutObject(ctx, minioStore.config.BookCoverBucket,
		path, contentBuffer, int64(contentBuffer.Len()), minio.PutObjectOptions{})

	return err
}
