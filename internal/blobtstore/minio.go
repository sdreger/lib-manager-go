package blobtstore

import (
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sdreger/lib-manager-go/internal/config"
	"io"
	"log/slog"
	"net/http"
)

type MinioStore struct {
	config          config.BLOBStoreConfig
	client          *minio.Client
	logger          *slog.Logger
	coverBucketName string
}

func NewMinioStore(logger *slog.Logger, config config.BLOBStoreConfig) (*MinioStore, error) {
	client, clientErr := getMinioClient(
		config.MinioEndpoint,
		config.MinioAccessKeyID,
		config.MinioSecretAccessKey,
		config.MinioUseSSL,
	)
	if clientErr != nil {
		return nil, clientErr
	}

	return &MinioStore{
		config:          config,
		client:          client,
		logger:          logger,
		coverBucketName: config.BookCoverBucket,
	}, nil
}

func (s *MinioStore) CreateBuckets(ctx context.Context) error {
	return createBookCoverBucketIfNotExist(ctx, s.logger, s.client, s.coverBucketName)
}

func (s *MinioStore) CoverExists(ctx context.Context, filePath string) bool {
	_, err := s.getObjectStats(ctx, s.coverBucketName, filePath)
	return err == nil || (errors.As(err, &minio.ErrorResponse{}) &&
		err.(minio.ErrorResponse).StatusCode != http.StatusNotFound)
}

func (s *MinioStore) GetBookCover(ctx context.Context, filePath string) (io.Reader, error) {
	return s.getObject(ctx, s.coverBucketName, filePath)
}

func (s *MinioStore) getObject(ctx context.Context, bucketName string, filePath string) (*minio.Object, error) {
	return s.client.GetObject(ctx, bucketName, filePath, minio.GetObjectOptions{})
}

func (s *MinioStore) getObjectStats(ctx context.Context, bucketName string, filePath string) (minio.ObjectInfo, error) {
	return s.client.StatObject(ctx, bucketName, filePath, minio.StatObjectOptions{})
}

func getMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*minio.Client, error) {
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
}

func createBookCoverBucketIfNotExist(ctx context.Context, logger *slog.Logger, client *minio.Client,
	bucketName string) error {

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !exists {
		logger.Info("creating book cover bucket", slog.String("bucket", bucketName))
		return client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}

	return nil
}
