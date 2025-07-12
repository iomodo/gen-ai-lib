package genailib

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

// UploadImageToS3 uploads image bytes to an Amazon S3 bucket and returns the public URL.
func UploadImageToS3(ctx context.Context, bucketName, objectName string, data []byte) (string, error) {
	return uploadToS3(ctx, bucketName, objectName, data)
}

// UploadVideoToS3 uploads video bytes to an Amazon S3 bucket and returns the public URL.
func UploadVideoToS3(ctx context.Context, bucketName, objectName string, data []byte) (string, error) {
	return uploadToS3(ctx, bucketName, objectName, data)
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
