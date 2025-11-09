package s3

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client *s3.Client
	bucket string
}

//func NewS3Client(ctx context.Context, bucket string) (*S3Client, error) {
//	cfg, err := config.LoadDefaultConfig(ctx)
//	if err != nil {
//		panic("unable to load SDK config, " + err.Error())
//	}
//
//	client := s3.NewFromConfig(cfg)
//	return &S3Client{client: client, bucket: bucket}, nil
//}

func (s *S3Client) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, path string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	uploader := manager.NewUploader(s.client)

	key := fmt.Sprintf("%s/%d_%s", path, time.Now().UnixNano(), fileHeader.Filename)

	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   file,
		ACL:    "public-read",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return result.Location, nil
}
