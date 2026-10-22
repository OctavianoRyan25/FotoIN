package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Photo struct {
	ID               uuid.UUID  `json:"id"`
	GalleryID        *uuid.UUID `json:"gallery_id,omitempty"`
	EventID          *uuid.UUID `json:"event_id,omitempty"`
	S3Key            string     `json:"-"`
	S3KeyCompressed  string     `json:"-"`
	S3KeyWatermarked string     `json:"-"`
	Filename         string     `json:"filename"`
	SizeBytes        int64      `json:"size_bytes"`
	Width            int        `json:"width"`
	Height           int        `json:"height"`
	FaceVectorID     *string    `json:"face_vector_id,omitempty"`
	PreviewURL       string     `json:"preview_url,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	Gallery *Gallery `json:"gallery,omitempty"`
	Event   *Event   `json:"event,omitempty"`
}

type PhotoRepository interface {
	Create(ctx context.Context, photo *Photo) error
	FindByID(ctx context.Context, id uuid.UUID) (*Photo, error)
	FindByGalleryID(ctx context.Context, galleryID uuid.UUID, page, limit int) ([]Photo, int64, error)
	FindByEventID(ctx context.Context, eventID uuid.UUID, page, limit int) ([]Photo, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]Photo, error)
	BulkDelete(ctx context.Context, galleryID uuid.UUID) error
}
