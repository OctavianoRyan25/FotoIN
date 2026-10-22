package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Gallery struct {
	ID             uuid.UUID  `json:"id"`
	PhotographerID uuid.UUID  `json:"photographer_id"`
	EventID        *uuid.UUID `json:"event_id,omitempty"`
	Title          string     `json:"title"`
	Slug           string     `json:"slug"`
	EventDate      time.Time  `json:"event_date"`
	Location       string     `json:"location"`
	Description    string     `json:"description"`
	PricePerPhoto  float64    `json:"price_per_photo"`
	CoverPhotoURL  string     `json:"cover_photo_url"`
	IsActive       bool       `json:"is_active"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	PhotoCount     int        `json:"photo_count"`
	TotalSales     int        `json:"total_sales,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Photographer *User   `json:"photographer,omitempty"`
	Event        *Event  `json:"event,omitempty"`
	Photos       []Photo `json:"photos,omitempty"`
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
