package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID             uuid.UUID  `json:"id"`
	PhotographerID uuid.UUID  `json:"photographer_id"`
	Title          string     `json:"title"`
	Slug           string     `json:"slug"`
	EventDate      *time.Time `json:"event_date,omitempty"`
	Location       string     `json:"location,omitempty"`
	Description    string     `json:"description,omitempty"`
	PricePerPhoto  float64    `json:"price_per_photo"`
	CoverPhotoURL  string     `json:"cover_photo_url,omitempty"`
	IsActive       bool       `json:"is_active"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Photographer *User     `json:"photographer,omitempty"`
	Galleries    []Gallery `json:"galleries,omitempty"`
	Photos       []Photo   `json:"photos,omitempty"`
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
