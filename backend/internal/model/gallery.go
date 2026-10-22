package model

import (
	"time"

	"fotoin/backend/internal/domain"

	"github.com/google/uuid"
)

type GalleryModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PhotographerID uuid.UUID  `gorm:"type:uuid;not null;index"`
	EventID       *uuid.UUID `gorm:"type:uuid;null;index"`
	Title         string     `gorm:"size:255;not null"`
	Slug          string     `gorm:"size:255;uniqueIndex;not null"`
	EventDate     time.Time  `gorm:"not null"`
	Location      string     `gorm:"size:255"`
	Description   string     `gorm:"type:text"`
	PricePerPhoto float64    `gorm:"type:decimal(12,2);not null"`
	CoverPhotoURL string     `gorm:"size:500"`
	IsActive      bool       `gorm:"default:true"`
	ExpiresAt     *time.Time `gorm:"null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Relations
	Photographer UserModel    `gorm:"foreignKey:PhotographerID"`
	Event        *EventModel  `gorm:"foreignKey:EventID"`
	Photos       []PhotoModel `gorm:"foreignKey:GalleryID"`
}

func (GalleryModel) TableName() string {
	return "galleries"
}

func (m *GalleryModel) ToDomain() *domain.Gallery {
	if m == nil {
		return nil
	}
	g := &domain.Gallery{
		ID:             m.ID,
		PhotographerID: m.PhotographerID,
		EventID:        m.EventID,
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
		g.Photographer = m.Photographer.ToDomain()
	}

	if m.Event != nil && m.Event.ID != uuid.Nil {
		g.Event = m.Event.ToDomain()
	}

	return g
}

func FromDomainGallery(g *domain.Gallery) *GalleryModel {
	if g == nil {
		return nil
	}
	return &GalleryModel{
		ID:             g.ID,
		PhotographerID: g.PhotographerID,
		EventID:        g.EventID,
		Title:          g.Title,
		Slug:           g.Slug,
		EventDate:      g.EventDate,
		Location:       g.Location,
		Description:    g.Description,
		PricePerPhoto:  g.PricePerPhoto,
		CoverPhotoURL:  g.CoverPhotoURL,
		IsActive:       g.IsActive,
		ExpiresAt:      g.ExpiresAt,
		CreatedAt:      g.CreatedAt,
		UpdatedAt:      g.UpdatedAt,
	}
}

func ToDomainGalleryList(models []GalleryModel) []domain.Gallery {
	res := make([]domain.Gallery, len(models))
	for i, m := range models {
		res[i] = *m.ToDomain()
	}
	return res
}
