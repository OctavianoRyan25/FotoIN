package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PhotographerID uuid.UUID  `gorm:"type:uuid;not null;index" json:"photographer_id"`
	Title          string     `gorm:"size:255;not null" json:"title"`
	Slug           string     `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	EventDate      *time.Time `gorm:"null" json:"event_date,omitempty"`
	Location       string     `gorm:"size:255" json:"location,omitempty"`
	Description    string     `gorm:"type:text" json:"description,omitempty"`
	PricePerPhoto  float64    `gorm:"type:decimal(12,2);not null;default:25000" json:"price_per_photo"`
	CoverPhotoURL  string     `gorm:"size:500" json:"cover_photo_url,omitempty"`
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	ExpiresAt      *time.Time `gorm:"null" json:"expires_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Photographer User      `gorm:"foreignKey:PhotographerID" json:"photographer,omitempty"`
	Galleries    []Gallery `gorm:"foreignKey:EventID" json:"galleries,omitempty"`
	Photos       []Photo   `gorm:"foreignKey:EventID" json:"photos,omitempty"`
}

type EventRepository interface {
	Create(ctx context.Context, event *Event) error
	FindByID(ctx context.Context, id uuid.UUID) (*Event, error)
	FindBySlug(ctx context.Context, slug string) (*Event, error)
	Update(ctx context.Context, event *Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByPhotographer(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]Event, int64, error)
	ListPublic(ctx context.Context, page, limit int, search string) ([]Event, int64, error)
}
