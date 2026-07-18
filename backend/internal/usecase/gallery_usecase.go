package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"

	"github.com/google/uuid"
)

type GalleryUsecase struct {
	galleryRepo domain.GalleryRepository
	userRepo    domain.UserRepository
	s3          *infra.S3Client
}

func NewGalleryUsecase(galleryRepo domain.GalleryRepository, userRepo domain.UserRepository, s3 *infra.S3Client) *GalleryUsecase {
	return &GalleryUsecase{galleryRepo: galleryRepo, userRepo: userRepo, s3: s3}
}

type CreateGalleryInput struct {
	Title         string    `json:"title" validate:"required,min=3,max=255"`
	EventDate     time.Time `json:"event_date" validate:"required"`
	Location      string    `json:"location" validate:"required"`
	Description   string    `json:"description"`
	PricePerPhoto float64   `json:"price_per_photo" validate:"required,min=10000"` // Floor price enforced here
}

type UpdateGalleryInput struct {
	Title         string    `json:"title" validate:"omitempty,min=3,max=255"`
	EventDate     time.Time `json:"event_date"`
	Location      string    `json:"location"`
	Description   string    `json:"description"`
	PricePerPhoto float64   `json:"price_per_photo" validate:"omitempty,min=10000"` // Floor price enforced here
	IsActive      *bool     `json:"is_active"`
}

func (u *GalleryUsecase) Create(ctx context.Context, photographerID uuid.UUID, input CreateGalleryInput) (*domain.Gallery, error) {
	// CRITICAL: Enforce floor price at usecase layer
	if input.PricePerPhoto < domain.MinPricePerPhoto {
		return nil, fmt.Errorf("price per photo must be at least Rp %.0f", domain.MinPricePerPhoto)
	}

	photographer, err := u.userRepo.FindByID(ctx, photographerID)
	if err != nil {
		return nil, errors.New("photographer not found")
	}

	slug := gallerySlug(input.Title, photographerID.String()[:8])

	gallery := &domain.Gallery{
		ID:             uuid.New(),
		PhotographerID: photographerID,
		Title:          input.Title,
		Slug:           slug,
		EventDate:      input.EventDate,
		Location:       input.Location,
		Description:    input.Description,
		PricePerPhoto:  input.PricePerPhoto,
		IsActive:       true,
	}

	// Set expiry for Basic tier
	if photographer.SubscriptionTier == domain.TierBasic {
		expiry := time.Now().Add(time.Duration(domain.BasicRetentionDays) * 24 * time.Hour)
		gallery.ExpiresAt = &expiry
	}

	if err := u.galleryRepo.Create(ctx, gallery); err != nil {
		return nil, err
	}

	return gallery, nil
}

func (u *GalleryUsecase) Update(ctx context.Context, id, photographerID uuid.UUID, input UpdateGalleryInput) (*domain.Gallery, error) {
	gallery, err := u.galleryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("gallery not found")
	}

	if gallery.PhotographerID != photographerID {
		return nil, errors.New("unauthorized: you don't own this gallery")
	}

	// CRITICAL: Floor price check
	if input.PricePerPhoto > 0 && input.PricePerPhoto < domain.MinPricePerPhoto {
		return nil, fmt.Errorf("price per photo must be at least Rp %.0f", domain.MinPricePerPhoto)
	}

	if input.Title != "" {
		gallery.Title = input.Title
	}
	if !input.EventDate.IsZero() {
		gallery.EventDate = input.EventDate
	}
	if input.Location != "" {
		gallery.Location = input.Location
	}
	if input.Description != "" {
		gallery.Description = input.Description
	}
	if input.PricePerPhoto > 0 {
		gallery.PricePerPhoto = input.PricePerPhoto
	}
	if input.IsActive != nil {
		gallery.IsActive = *input.IsActive
	}

	if err := u.galleryRepo.Update(ctx, gallery); err != nil {
		return nil, err
	}

	return gallery, nil
}

func (u *GalleryUsecase) Delete(ctx context.Context, id, photographerID uuid.UUID) error {
	gallery, err := u.galleryRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("gallery not found")
	}
	if gallery.PhotographerID != photographerID {
		return errors.New("unauthorized: you don't own this gallery")
	}
	return u.galleryRepo.Delete(ctx, id)
}

func (u *GalleryUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Gallery, error) {
	return u.galleryRepo.FindByID(ctx, id)
}

func (u *GalleryUsecase) GetBySlug(ctx context.Context, slug string, s3 *infra.S3Client) (*domain.Gallery, error) {
	gallery, err := u.galleryRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("gallery not found or inactive")
	}

	// Check expiry
	if gallery.ExpiresAt != nil && gallery.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("this gallery has expired")
	}

	return gallery, nil
}

func (u *GalleryUsecase) ListByPhotographer(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]domain.Gallery, int64, error) {
	return u.galleryRepo.FindByPhotographerID(ctx, photographerID, page, limit)
}

func (u *GalleryUsecase) ListPublic(ctx context.Context, page, limit int) ([]domain.Gallery, int64, error) {
	return u.galleryRepo.ListPublic(ctx, page, limit)
}

func gallerySlug(title, suffix string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	slug = result.String()
	slug = strings.Trim(slug, "-")
	return fmt.Sprintf("%s-%s", slug, suffix)
}
