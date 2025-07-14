package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Storage interface {
	Upload(ctx context.Context, data []byte, objectName string) (string, error)
}

// generateObjectName returns objectName if provided, otherwise a random UUID-based name.
func generateObjectName(objectName string) string {
	if objectName != "" {
		return objectName
	}
	return fmt.Sprintf("object-%s", uuid.New().String())
}
