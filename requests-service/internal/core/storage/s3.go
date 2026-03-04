package storage

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path"

	"github.com/RadekKusiak71/splitwise-requests/internal/core/config"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(region string) (*s3.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(region))
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	return client, nil
}

type S3PutObjectAPI interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type S3Uploader struct {
	s3     S3PutObjectAPI
	config *config.AWSConfig
}

func NewS3Uploader(s3client S3PutObjectAPI, config *config.AWSConfig) *S3Uploader {
	return &S3Uploader{s3: s3client, config: config}
}

func (u *S3Uploader) Upload(ctx context.Context, key string, file io.Reader) error {
	ext := path.Ext(key)
	contentType := mime.TypeByExtension(ext)

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := u.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &u.config.BucketName,
		Key:         &key,
		Body:        file,
		ContentType: &contentType,
	})

	if err != nil {
		return fmt.Errorf("s3 put object bucket=%s key=%s: %w", u.config.BucketName, key, err)
	}
	return nil
}

func (u *S3Uploader) GenerateObjectURL(key string) string {
	if key != "" {
		return fmt.Sprintf("%s%s", u.config.S3BaseURL, key)
	}
	return ""
}
