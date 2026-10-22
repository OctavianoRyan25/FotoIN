package domain

import (
	"io"
	"time"
)

// StorageService defines the contract for object storage operations (e.g. S3 / MinIO).
type StorageService interface {
	GetPresignedUploadURL(key, contentType string, ttl time.Duration) (string, error)
	GetPresignedDownloadURL(key string, ttl time.Duration) (string, error)
	UploadObject(key, contentType string, body io.Reader, size int64) error
	DeleteObject(key string) error
	GetPublicURL(key string) string
}
