package model

import (
	"time"

	"fotoin/backend/internal/domain"

	"github.com/google/uuid"
)

type PhotoModel struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	GalleryID        *uuid.UUID `gorm:"type:uuid;null;index"`
	EventID          *uuid.UUID `gorm:"type:uuid;null;index"`
	S3Key            string     `gorm:"size:500;not null"`
	S3KeyCompressed  string     `gorm:"size:500"`
	S3KeyWatermarked string     `gorm:"size:500"`
	Filename         string     `gorm:"size:255"`
	SizeBytes        int64      `gorm:"default:0"`
	Width            int        `gorm:"default:0"`
	Height           int        `gorm:"default:0"`
	FaceVectorID     *string    `gorm:"size:255;null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Relations
	Gallery *GalleryModel `gorm:"foreignKey:GalleryID"`
	Event   *EventModel   `gorm:"foreignKey:EventID"`
}

func (PhotoModel) TableName() string {
	return "photos"
}

func (m *PhotoModel) ToDomain() *domain.Photo {
	if m == nil {
		return nil
	}
	p := &domain.Photo{
		ID:               m.ID,
		GalleryID:        m.GalleryID,
		EventID:          m.EventID,
		S3Key:            m.S3Key,
		S3KeyCompressed:  m.S3KeyCompressed,
		S3KeyWatermarked: m.S3KeyWatermarked,
		Filename:         m.Filename,
		SizeBytes:        m.SizeBytes,
		Width:            m.Width,
		Height:           m.Height,
		FaceVectorID:     m.FaceVectorID,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}

	if m.Gallery != nil && m.Gallery.ID != uuid.Nil {
		p.Gallery = m.Gallery.ToDomain()
	}
	if m.Event != nil && m.Event.ID != uuid.Nil {
		p.Event = m.Event.ToDomain()
	}

	return p
}

func FromDomainPhoto(p *domain.Photo) *PhotoModel {
	if p == nil {
		return nil
	}
	return &PhotoModel{
		ID:               p.ID,
		GalleryID:        p.GalleryID,
		EventID:          p.EventID,
		S3Key:            p.S3Key,
		S3KeyCompressed:  p.S3KeyCompressed,
		S3KeyWatermarked: p.S3KeyWatermarked,
		Filename:         p.Filename,
		SizeBytes:        p.SizeBytes,
		Width:            p.Width,
		Height:           p.Height,
		FaceVectorID:     p.FaceVectorID,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
}

func ToDomainPhotoList(models []PhotoModel) []domain.Photo {
	res := make([]domain.Photo, len(models))
	for i, m := range models {
		res[i] = *m.ToDomain()
	}
	return res
}
