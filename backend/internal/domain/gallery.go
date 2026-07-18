package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Gallery struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PhotographerID uuid.UUID  `gorm:"type:uuid;not null;index" json:"photographer_id"`
	EventID        *uuid.UUID `gorm:"type:uuid;null;index" json:"event_id,omitempty"`
	Title          string     `gorm:"size:255;not null" json:"title"`
	Slug           string     `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	EventDate      time.Time  `gorm:"not null" json:"event_date"`
	Location       string     `gorm:"size:255" json:"location"`
	Description    string     `gorm:"type:text" json:"description"`
	PricePerPhoto  float64    `gorm:"type:decimal(12,2);not null" json:"price_per_photo"`
	CoverPhotoURL  string     `gorm:"size:500" json:"cover_photo_url"`
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	ExpiresAt      *time.Time `gorm:"null" json:"expires_at,omitempty"`
	PhotoCount     int        `gorm:"-" json:"photo_count"`
	TotalSales     int        `gorm:"-" json:"total_sales,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Photographer User    `gorm:"foreignKey:PhotographerID" json:"photographer,omitempty"`
	Event        *Event  `gorm:"foreignKey:EventID" json:"event,omitempty"`
	Photos       []Photo `gorm:"foreignKey:GalleryID" json:"-"`
}

type GalleryRepository interface {
	Create(ctx context.Context, gallery *Gallery) error
	FindByID(ctx context.Context, id uuid.UUID) (*Gallery, error)
	FindBySlug(ctx context.Context, slug string) (*Gallery, error)
	FindByPhotographerID(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]Gallery, int64, error)
	Update(ctx context.Context, gallery *Gallery) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListPublic(ctx context.Context, page, limit int) ([]Gallery, int64, error)
	CountPhotos(ctx context.Context, galleryID uuid.UUID) (int, error)
}
