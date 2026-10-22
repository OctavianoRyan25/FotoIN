package model

import (
	"time"

	"fotoin/backend/internal/domain"

	"github.com/google/uuid"
)

type EventModel struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PhotographerID uuid.UUID  `gorm:"type:uuid;not null;index"`
	Title          string     `gorm:"size:255;not null"`
	Slug           string     `gorm:"size:255;uniqueIndex;not null"`
	EventDate      *time.Time `gorm:"null"`
	Location       string     `gorm:"size:255"`
	Description    string     `gorm:"type:text"`
	PricePerPhoto  float64    `gorm:"type:decimal(12,2);not null;default:25000"`
	CoverPhotoURL  string     `gorm:"size:500"`
	IsActive       bool       `gorm:"default:true"`
	ExpiresAt      *time.Time `gorm:"null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// Relations
	Photographer UserModel      `gorm:"foreignKey:PhotographerID"`
	Galleries    []GalleryModel `gorm:"foreignKey:EventID"`
	Photos       []PhotoModel   `gorm:"foreignKey:EventID"`
}

func (EventModel) TableName() string {
	return "events"
}

func (m *EventModel) ToDomain() *domain.Event {
	if m == nil {
		return nil
	}
	ev := &domain.Event{
		ID:             m.ID,
		PhotographerID: m.PhotographerID,
		Title:          m.Title,
		Slug:           m.Slug,
		EventDate:      m.EventDate,
		Location:       m.Location,
		Description:    m.Description,
		PricePerPhoto:  m.PricePerPhoto,
		CoverPhotoURL:  m.CoverPhotoURL,
		IsActive:       m.IsActive,
		ExpiresAt:      m.ExpiresAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}

	if m.Photographer.ID != uuid.Nil {
		ev.Photographer = m.Photographer.ToDomain()
	}

	return ev
}

func FromDomainEvent(e *domain.Event) *EventModel {
	if e == nil {
		return nil
	}
	return &EventModel{
		ID:             e.ID,
		PhotographerID: e.PhotographerID,
		Title:          e.Title,
		Slug:           e.Slug,
		EventDate:      e.EventDate,
		Location:       e.Location,
		Description:    e.Description,
		PricePerPhoto:  e.PricePerPhoto,
		CoverPhotoURL:  e.CoverPhotoURL,
		IsActive:       e.IsActive,
		ExpiresAt:      e.ExpiresAt,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

func ToDomainEventList(models []EventModel) []domain.Event {
	res := make([]domain.Event, len(models))
	for i, m := range models {
		res[i] = *m.ToDomain()
	}
	return res
}
