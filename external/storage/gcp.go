package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type gcpService struct {
	client     *storage.Client
	bucketName string
}

func NewGCPService(bucketName string) (Storage, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create storage client")
	}
	return &gcpService{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (g *gcpService) Upload(ctx context.Context, data []byte, objectName string) (string, error) {
	var objName string
	if objectName != "" {
		objName = objectName
	} else {
		objName = fmt.Sprintf("object-%s", uuid.New().String())
	}
	return uploadToGCP(ctx, g.bucketName, objName, data)
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
