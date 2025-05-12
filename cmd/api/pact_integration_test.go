//go:build integration

package main

import (
	"fmt"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/pact-foundation/pact-go/v2/utils"
	"github.com/sdreger/lib-manager-go/internal/blobtstore"
	"github.com/sdreger/lib-manager-go/internal/config"
	"github.com/sdreger/lib-manager-go/internal/database"
	"github.com/sdreger/lib-manager-go/internal/tests"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	// max time the framework will wait to issue requests to your provider API
	pactRequestTimeout               = 5 * time.Second
	applicationStartupTimeoutSeconds = 5
	freeLocalHost                    = "127.0.0.1"
)

var (
	freeLocalPort, _ = utils.GetFreePort()
)

func TestPactProvider(t *testing.T) {

	// Start all the dependencies (Postgres/Mino), seed test data,
	// initialize and start the application instance
	startInstrumentedProvider(t)

	// Verify the Provider - Branch-based Published Pacts for any known consumers
	verifier := provider.NewVerifier()
	request := preparePactVerifyRequest()
	err := verifier.VerifyProvider(t, request)
	if err != nil {
		t.Fail()
	}
}

func startInstrumentedProvider(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// Initialize local HTTP server
	_ = os.Setenv(getEnvKey("HTTP_HOST"), freeLocalHost)
	_ = os.Setenv(getEnvKey("HTTP_PORT"), strconv.Itoa(freeLocalPort))

	// Initialize Postgres test container
	postgresContainer := tests.StartDBTestContainer(t)
	dbConfig := tests.GetTestDBConfig(t, postgresContainer)
	_ = os.Setenv(getEnvKey("DB_HOST"), dbConfig.Host)
	_ = os.Setenv(getEnvKey("DB_PORT"), strconv.Itoa(dbConfig.Port))
	_ = os.Setenv(getEnvKey("DB_USER"), dbConfig.User)
	_ = os.Setenv(getEnvKey("DB_PASSWORD"), dbConfig.Password)
	_ = os.Setenv(getEnvKey("DB_NAME"), dbConfig.Name)
	_ = os.Setenv(getEnvKey("DB_AUTO_MIGRATE"), "true")
	connection, err := database.Open(dbConfig)
	require.NoError(t, err, "failed to open database connection")
	err = database.Migrate(logger, dbConfig, connection.DB)
	require.NoError(t, err, "failed to perform database migration")

	// Fill the database with test data
	err = prepareTestData(postgresContainer, "testdata/book_lookup_filter.sql")
	require.NoError(t, err, "failed to load test SQL file")

	// Initialize Minio test container
	minioContainer := tests.StartMinioTestContainer(t)
	minioConfig := tests.GetTestMinioConfig(t, minioContainer)
	_ = os.Setenv(getEnvKey("BLOB_STORE_MINIO_ENDPOINT"), minioConfig.MinioEndpoint)
	_ = os.Setenv(getEnvKey("BLOB_STORE_MINIO_ACCESS_KEY_ID"), minioConfig.MinioSecretAccessKey)
	_ = os.Setenv(getEnvKey("BLOB_STORE_MINIO_ACCESS_SECRET_KEY"), minioConfig.MinioSecretAccessKey)
	minioStore, err := blobtstore.NewMinioStore(logger, minioConfig)
	require.NoError(t, err, "failed to create minio store")
	defer minioStore.Close()

	// Start the application instance asynchronously
	go func() {
		defer func() {
			_ = os.Unsetenv(getEnvKey("HTTP_HOST"))
			_ = os.Unsetenv(getEnvKey("HTTP_PORT"))
			_ = os.Unsetenv(getEnvKey("DB_HOST"))
			_ = os.Unsetenv(getEnvKey("DB_PORT"))
			_ = os.Unsetenv(getEnvKey("DB_USER"))
			_ = os.Unsetenv(getEnvKey("DB_PASSWORD"))
			_ = os.Unsetenv(getEnvKey("DB_NAME"))
			_ = os.Unsetenv(getEnvKey("DB_AUTO_MIGRATE"))
			_ = os.Unsetenv(getEnvKey("BLOB_STORE_MINIO_ENDPOINT"))
			_ = os.Unsetenv(getEnvKey("BLOB_STORE_MINIO_ACCESS_KEY_ID"))
			_ = os.Unsetenv(getEnvKey("BLOB_STORE_MINIO_ACCESS_SECRET_KEY"))
		}()
		err := run(logger)
		require.NoError(t, err)
	}()

	// Wait until the application and all dependencies to be ready
	waitForReadinessCheckPass(t, logger)
}

func getEnvKey(key string) string {
	return config.EnvPrefix + key
}

func prepareTestData(testContainer *postgres.PostgresContainer, fileName string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	return tests.ExecSQL(testContainer, string(file))
}

func preparePactVerifyRequest() provider.VerifyRequest {

	request := provider.VerifyRequest{
		ProviderBaseURL:            fmt.Sprintf("http://%s:%d", freeLocalHost, freeLocalPort),
		ProviderBranch:             os.Getenv("PACT_VERSION_BRANCH"),
		Provider:                   os.Getenv("PACT_PROVIDER_NAME"),
		FailIfNoPactsFound:         false,
		PublishVerificationResults: shouldPublishResults(),
		ProviderVersion:            os.Getenv("PACT_VERSION_COMMIT"),
		EnablePending:              true,
		RequestTimeout:             pactRequestTimeout,
		DisableColoredOutput:       true,
	}

	brokerUsername := os.Getenv("PACT_BROKER_USERNAME")
	if brokerUsername != "" {
		request.BrokerUsername = brokerUsername
	}

	brokerPassword := os.Getenv("PACT_BROKER_PASSWORD")
	if brokerPassword != "" {
		request.BrokerPassword = brokerPassword
	}

	pactURL := os.Getenv("PACT_URL")
	if pactURL != "" {
		request.PactURLs = append(request.PactURLs, pactURL)
	} else {
		request.BrokerURL = os.Getenv("PACT_BROKER_URL")
	}

	return request
}

func shouldPublishResults() bool {
	publishResults := os.Getenv("PACT_PUBLISH_RESULTS")
	if strings.ToLower(publishResults) == "true" {
		return true
	}

	return false
}

func waitForReadinessCheckPass(t *testing.T, logger *slog.Logger) {
	httpClient := &http.Client{}
	url := fmt.Sprintf("http://%s:%d/readyz", freeLocalHost, freeLocalPort)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < applicationStartupTimeoutSeconds; i++ {
		response, err := httpClient.Do(request)
		if err == nil && response.StatusCode == http.StatusOK {
			return
		}
		logger.Info("waiting for readiness check...", "url", url)
		time.Sleep(1 * time.Second)
	}

	t.Fatalf("failed to connect to http://%s:%d", freeLocalHost, freeLocalPort)
}
