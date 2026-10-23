package domain

import (
	"io"
	"time"
)

// StorageService defines the contract for object storage operations (e.g. AWS S3).
type StorageService interface {
	GetPresignedUploadURL(key string, contentType string, ttl time.Duration) (string, error)
	GetPresignedDownloadURL(key string, ttl time.Duration) (string, error)
	UploadObject(key string, contentType string, body io.Reader, size int64) error
	DeleteObject(key string) error
	GetPublicURL(key string) string
}

