package storage

import "context"

type Storage interface {
	Upload(ctx context.Context, data []byte, objectName string) (string, error)
}
