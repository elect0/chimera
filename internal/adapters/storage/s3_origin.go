package storage

import (
	"context"
	"io"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/elect0/chimera/internal/config"
	"github.com/elect0/chimera/internal/ports"
)

type S3OriginRepository struct {
	s3Client   *s3.Client
	bucketName string
	log        *slog.Logger
}

func NewS3OriginRepository(ctx context.Context, cfg *config.Config, log *slog.Logger) (*S3OriginRepository, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.S3.Region))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg)

	return &S3OriginRepository{
		s3Client:   s3Client,
		bucketName: cfg.S3.Bucket,
		log:        log,
	}, nil
}

func (r *S3OriginRepository) Get(ctx context.Context, imagePath string) ([]byte, error) {
	log := r.log.With(slog.String("imagePath", imagePath), slog.String("bucket", r.bucketName))
	log.Debug("fetching image from s3")

	result, err := r.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(imagePath),
	})
	if err != nil {
		log.Error("failed to get object from s3", slog.String("error", err.Error()))
		return nil, err
	}

	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Error("failed to read object body", slog.String("error", err.Error()))
		return nil, err
	}

	log.Debug("successfully fetched image from s3", slog.Int("size_bytes", len(body)))
	return body, nil
}

var _ ports.OriginRepository = (*S3OriginRepository)(nil)
