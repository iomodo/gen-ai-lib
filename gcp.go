package genailib

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

// UploadImageToGCP uploads image bytes to a GCP Cloud Storage bucket.
func UploadImageToGCP(ctx context.Context, bucketName, objectName string, data []byte) (string, error) {
	return uploadToGCP(ctx, bucketName, objectName, data)
}

// UploadVideoToGCP uploads video bytes to a GCP Cloud Storage bucket.
func UploadVideoToGCP(ctx context.Context, bucketName, objectName string, data []byte) (string, error) {
	return uploadToGCP(ctx, bucketName, objectName, data)
}

func uploadToGCP(ctx context.Context, bucketName, objectName string, data []byte) (string, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to create storage client")
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
		w.Close()
		return "", errors.Wrap(err, "failed to write data")
	}
	if err := w.Close(); err != nil {
		return "", errors.Wrap(err, "failed to close writer")
	}
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", errors.Wrap(err, "failed to set object ACL")
	}
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)
	return url, nil
}
