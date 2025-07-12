package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type s3Service struct {
	client     *s3.Client
	bucketName string
}

func NewS3Service(bucketName string) (Storage, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to load AWS config")
	}
	return &s3Service{
		client:     s3.NewFromConfig(cfg),
		bucketName: bucketName,
	}, nil
}

func (s *s3Service) Upload(ctx context.Context, data []byte, objectName string) (string, error) {
	var objName string
	if objectName != "" {
		objName = objectName
	} else {
		objName = fmt.Sprintf("object-%s", uuid.New().String())
	}
	return uploadToS3(ctx, s.bucketName, objName, data)
}

func uploadToS3(ctx context.Context, bucketName, objectName string, data []byte) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to load AWS config")
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
		Body:   bytes.NewReader(data),
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to put object")
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, objectName)
	return url, nil
}
