package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Photo struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	GalleryID        *uuid.UUID `gorm:"type:uuid;null;index" json:"gallery_id,omitempty"`
	EventID          *uuid.UUID `gorm:"type:uuid;null;index" json:"event_id,omitempty"`
	S3Key            string     `gorm:"size:500;not null" json:"-"`     // Original high-res (private)
	S3KeyCompressed  string     `gorm:"size:500" json:"-"`              // Compressed ~2-5MB (for delivery)
	S3KeyWatermarked string     `gorm:"size:500" json:"-"`              // Watermarked (for preview)
	Filename         string     `gorm:"size:255" json:"filename"`
	SizeBytes        int64      `gorm:"default:0" json:"size_bytes"`
	Width            int        `gorm:"default:0" json:"width"`
	Height           int        `gorm:"default:0" json:"height"`
	FaceVectorID     *string    `gorm:"size:255;null" json:"face_vector_id,omitempty"` // Fase 2: ML
	PreviewURL       string     `gorm:"-" json:"preview_url,omitempty"`               // Signed URL (ephemeral)
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Relations
	Gallery *Gallery `gorm:"foreignKey:GalleryID" json:"-"`
	Event   *Event   `gorm:"foreignKey:EventID" json:"-"`
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
