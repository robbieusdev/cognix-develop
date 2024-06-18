package storage

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

type (
	MinioConfig struct {
		AccessKey       string `env:"MINIO_ACCESS_KEY"`
		SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY"`
		Endpoint        string `env:"MINIO_ENDPOINT"`
		UseSSL          bool   `env:"MINIO_USE_SSL"`
		BucketName      string `env:"MINIO_BUCKET_NAME"`
		Region          string `env:"MINIO_REGION"`
	}
	MinIOClient interface {
		Upload(ctx context.Context, bucket, filename, contentType string, reader io.Reader) (string, string, error)
		GetObject(ctx context.Context, bucket, filename string, writer io.Writer) error
	}
	minIOClient struct {
		Region string
		client *minio.Client
	}
	minIOMockClient struct{}
)

func (c *minIOClient) checkOrCreateBucket(ctx context.Context, bucketName string) error {
	ok, err := c.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	// create bucket if not exists
	if !ok {
		return c.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: c.Region,
		})
	}
	return nil
}
func (c *minIOClient) Upload(ctx context.Context, bucket, filename, contentType string, reader io.Reader) (string, string, error) {
	// verify is bucket exists. create if not exists
	if err := c.checkOrCreateBucket(ctx, bucket); err != nil {
		return "", "", err
	}

	// save file in minio
	res, err := c.client.PutObject(ctx, bucket, filename, reader, -1,
		minio.PutObjectOptions{
			ContentType: contentType,
			NumThreads:  0,
		})
	if err != nil {
		return "", "", utils.Internal.Wrapf(err, "cannot upload file: %s", err.Error())
	}
	return res.Key, res.ChecksumCRC32C, nil
}

func (c *minIOClient) GetObject(ctx context.Context, bucket, filename string, writer io.Writer) error {
	object, err := c.client.GetObject(ctx, bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()
	_, err = io.Copy(writer, object)
	if err != nil {
		return err
	}
	return nil
}

func NewMinIOClient(cfg *MinioConfig) (MinIOClient, error) {

	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.AccessKey,
			cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &minIOClient{
		Region: cfg.Region,
		client: minioClient,
	}, nil
}
