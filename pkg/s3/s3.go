package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mpb/configs"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client struct {
	client   *s3.Client
	uploader *manager.Uploader
	bucket   string
	region   string
}

func NewS3Client(conf *configs.Config) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(conf.AWS.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)

	return &S3Client{
		client:   client,
		uploader: uploader,
		bucket:   conf.AWS.Bucket,
		region:   conf.AWS.Region,
	}, nil
}

func (s *S3Client) UploadFile(ctx context.Context, key string, file io.Reader, contentType string) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	_, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, url.PathEscape(key))
	return url, nil
}
