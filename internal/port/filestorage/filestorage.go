package filestorage

import (
	"context"
	"time"
)

//go:generate go tool mockgen -source=filestorage.go -destination=../../../mocks/mock_filestorage.go -package=mocks

// FileStorage defines the interface for file storage operations.
type FileStorage interface {
	// Upload uploads a file to the storage and returns the object name.
	Upload(ctx context.Context, data []byte, mimeType string) (string, error)

	// GetPreSignedURL generates a pre-signed URL for the given object.
	// The URL will be valid for the specified duration.
	GetPreSignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)

	// Delete removes a file from the storage.
	Delete(ctx context.Context, objectName string) error
}
