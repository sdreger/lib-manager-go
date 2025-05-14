//go:build !build

package tests

import (
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	testMinio "github.com/testcontainers/testcontainers-go/modules/minio"
	"log"
	"testing"
)

func StartMinioTestContainer(t *testing.T) *testMinio.MinioContainer {
	minioContainer, err := testMinio.Run(t.Context(), "quay.io/minio/minio:RELEASE.2025-02-18T16-25-55Z")
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(minioContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	})
	require.NoError(t, err, "failed to start minio container")

	return minioContainer
}

func GetTestMinioConfig(t *testing.T, minioContainer *testMinio.MinioContainer) config.BLOBStoreConfig {
	endpoint, err := minioContainer.ConnectionString(t.Context())
	require.NoError(t, err, "failed to get container endpoint")

	appConfig, err := config.New()
	require.NoError(t, err, "failed to load app config")
	blobStoreConfig := appConfig.BLOBStore
	blobStoreConfig.MinioEndpoint = endpoint
	blobStoreConfig.MinioAccessKeyID = minioContainer.Username
	blobStoreConfig.MinioSecretAccessKey = minioContainer.Password

	return blobStoreConfig
}

func TerminateMinioContainer(t *testing.T, minioContainer *testMinio.MinioContainer) {
	err := testcontainers.TerminateContainer(minioContainer)
	require.NoError(t, err, "failed to terminate container")
}
