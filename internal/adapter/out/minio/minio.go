package minio

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config holds MinIO configuration.
type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	PublicDomain    string // e.g., "s3.vladimish.com"
}

// FileStorage implements the filestorage.FileStorage interface using MinIO.
type FileStorage struct {
	client       *minio.Client
	bucketName   string
	publicDomain string
}

// NewFileStorage creates a new MinIO file storage instance.
func NewFileStorage(cfg Config) (*FileStorage, error) {
	// Initialize MinIO client
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	// Check if bucket exists, create if not
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &FileStorage{
		client:       client,
		bucketName:   cfg.BucketName,
		publicDomain: cfg.PublicDomain,
	}, nil
}

// Upload uploads a file to MinIO and returns the object name.
func (fs *FileStorage) Upload(ctx context.Context, data []byte, mimeType string) (string, error) {
	// Generate unique object name with extension based on MIME type
	objectName := generateObjectName(mimeType)

	// Upload the file
	reader := bytes.NewReader(data)
	_, err := fs.client.PutObject(ctx, fs.bucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: mimeType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	return objectName, nil
}

// GetPreSignedURL generates a pre-signed URL for the given object.
func (fs *FileStorage) GetPreSignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	// Generate pre-signed URL
	reqParams := make(url.Values)
	presignedURL, err := fs.client.PresignedGetObject(ctx, fs.bucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %w", err)
	}

	// If we have a public domain, replace the endpoint with it
	if fs.publicDomain != "" {
		u, parseErr := url.Parse(presignedURL.String())
		if parseErr != nil {
			return "", fmt.Errorf("failed to parse pre-signed URL: %w", parseErr)
		}
		u.Host = fs.publicDomain
		u.Scheme = "https" // Always use HTTPS for public domain
		return u.String(), nil
	}

	return presignedURL.String(), nil
}

// Delete removes a file from MinIO.
func (fs *FileStorage) Delete(ctx context.Context, objectName string) error {
	err := fs.client.RemoveObject(ctx, fs.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// generateObjectName generates a unique object name with appropriate extension.
func generateObjectName(mimeType string) string {
	id := uuid.New().String()

	// Map common MIME types to extensions
	var ext string
	switch mimeType {
	case "image/jpeg", "image/jpg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	case "image/bmp":
		ext = ".bmp"
	default:
		ext = ".bin"
	}

	// Create path with date prefix for better organization
	now := time.Now()
	return fmt.Sprintf("images/%d/%02d/%02d/%s%s", now.Year(), now.Month(), now.Day(), id, ext)
}
