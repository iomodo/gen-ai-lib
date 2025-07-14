package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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
	objName := generateObjectName(objectName)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objName),
		Body:   bytes.NewReader(data),
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to put object")
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, objName)
	return url, nil
}
