package usecase

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"

	"github.com/google/uuid"
)

type PhotoUsecase struct {
	photoRepo   domain.PhotoRepository
	galleryRepo domain.GalleryRepository
	eventRepo   domain.EventRepository
	userRepo    domain.UserRepository
	s3          *infra.S3Client
}

func NewPhotoUsecase(
	photoRepo domain.PhotoRepository,
	galleryRepo domain.GalleryRepository,
	eventRepo domain.EventRepository,
	userRepo domain.UserRepository,
	s3 *infra.S3Client,
) *PhotoUsecase {
	return &PhotoUsecase{
		photoRepo:   photoRepo,
		galleryRepo: galleryRepo,
		eventRepo:   eventRepo,
		userRepo:    userRepo,
		s3:          s3,
	}
}

// ─── Gallery Photos ────────────────────────────────────────────────────────────

type GetUploadURLInput struct {
	GalleryID   uuid.UUID `json:"gallery_id" validate:"required"`
	Filename    string    `json:"filename" validate:"required"`
	ContentType string    `json:"content_type" validate:"required"`
}

type UploadURLResponse struct {
	UploadURL    string `json:"upload_url"`
	S3Key        string `json:"s3_key"`
	S3KeyPreview string `json:"s3_key_preview"`
}

type ConfirmUploadInput struct {
	GalleryID string `json:"gallery_id" validate:"required"`
	S3Key     string `json:"s3_key" validate:"required"`
	Filename  string `json:"filename" validate:"required"`
	SizeBytes int64  `json:"size_bytes" validate:"required,min=1"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func (u *PhotoUsecase) GetUploadURL(ctx context.Context, photographerID uuid.UUID, input GetUploadURLInput) (*UploadURLResponse, error) {
	gallery, err := u.galleryRepo.FindByID(ctx, input.GalleryID)
	if err != nil {
		return nil, errors.New("gallery not found")
	}
	if gallery.PhotographerID != photographerID {
		return nil, errors.New("unauthorized: you don't own this gallery")
	}

	photographer, err := u.userRepo.FindByID(ctx, photographerID)
	if err != nil {
		return nil, err
	}
	if photographer.StorageUsedBytes >= photographer.StorageLimit() {
		return nil, fmt.Errorf("storage limit exceeded (%.1f GB used of %.0f GB)",
			photographer.StorageUsedGB(),
			float64(photographer.StorageLimit())/(1024*1024*1024))
	}

	photoID := uuid.New()
	ext := filepath.Ext(input.Filename)
	s3Key := fmt.Sprintf("galleries/%s/photos/%s%s", input.GalleryID, photoID.String(), ext)
	s3KeyPreview := fmt.Sprintf("galleries/%s/watermarked/%s%s", input.GalleryID, photoID.String(), ext)

	uploadURL, err := u.s3.GetPresignedUploadURL(s3Key, input.ContentType, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	return &UploadURLResponse{
		UploadURL:    uploadURL,
		S3Key:        s3Key,
		S3KeyPreview: s3KeyPreview,
	}, nil
}

func (u *PhotoUsecase) ConfirmUpload(ctx context.Context, photographerID uuid.UUID, input ConfirmUploadInput) (*domain.Photo, error) {
	galleryID, err := uuid.Parse(input.GalleryID)
	if err != nil {
		return nil, errors.New("invalid gallery ID")
	}
	gallery, err := u.galleryRepo.FindByID(ctx, galleryID)
	if err != nil {
		return nil, errors.New("gallery not found")
	}
	if gallery.PhotographerID != photographerID {
		return nil, errors.New("unauthorized")
	}

	s3KeyWatermarked := strings.Replace(input.S3Key, "/photos/", "/watermarked/", 1)

	photo := &domain.Photo{
		ID:               uuid.New(),
		GalleryID:        &galleryID,
		S3Key:            input.S3Key,
		S3KeyWatermarked: s3KeyWatermarked,
		Filename:         input.Filename,
		SizeBytes:        input.SizeBytes,
		Width:            input.Width,
		Height:           input.Height,
	}

	if err := u.photoRepo.Create(ctx, photo); err != nil {
		return nil, err
	}
	u.userRepo.UpdateStorageUsed(ctx, photographerID, input.SizeBytes)

	return photo, nil
}

func (u *PhotoUsecase) ListPhotos(ctx context.Context, galleryID uuid.UUID, page, limit int) ([]domain.Photo, int64, error) {
	photos, total, err := u.photoRepo.FindByGalleryID(ctx, galleryID, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return attachPreviewURLs(photos, u.s3), total, nil
}

// ─── Event Photos ─────────────────────────────────────────────────────────────

type GetEventUploadURLInput struct {
	EventID     uuid.UUID `json:"event_id" validate:"required"`
	Filename    string    `json:"filename" validate:"required"`
	ContentType string    `json:"content_type" validate:"required"`
}

type ConfirmEventUploadInput struct {
	EventID   string `json:"event_id" validate:"required"`
	S3Key     string `json:"s3_key" validate:"required"`
	Filename  string `json:"filename" validate:"required"`
	SizeBytes int64  `json:"size_bytes" validate:"required,min=1"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func (u *PhotoUsecase) GetEventUploadURL(ctx context.Context, photographerID uuid.UUID, input GetEventUploadURLInput) (*UploadURLResponse, error) {
	if u.eventRepo == nil {
		return nil, errors.New("event repository not available")
	}
	event, err := u.eventRepo.FindByID(ctx, input.EventID)
	if err != nil {
		return nil, errors.New("event not found")
	}
	if event.PhotographerID != photographerID {
		return nil, errors.New("unauthorized: you don't own this event")
	}

	photographer, err := u.userRepo.FindByID(ctx, photographerID)
	if err != nil {
		return nil, err
	}
	if photographer.StorageUsedBytes >= photographer.StorageLimit() {
		return nil, fmt.Errorf("storage limit exceeded")
	}

	photoID := uuid.New()
	ext := filepath.Ext(input.Filename)
	s3Key := fmt.Sprintf("events/%s/photos/%s%s", input.EventID, photoID.String(), ext)
	s3KeyPreview := fmt.Sprintf("events/%s/watermarked/%s%s", input.EventID, photoID.String(), ext)

	uploadURL, err := u.s3.GetPresignedUploadURL(s3Key, input.ContentType, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	return &UploadURLResponse{
		UploadURL:    uploadURL,
		S3Key:        s3Key,
		S3KeyPreview: s3KeyPreview,
	}, nil
}

func (u *PhotoUsecase) ConfirmEventUpload(ctx context.Context, photographerID uuid.UUID, input ConfirmEventUploadInput) (*domain.Photo, error) {
	eventID, err := uuid.Parse(input.EventID)
	if err != nil {
		return nil, errors.New("invalid event ID")
	}
	if u.eventRepo == nil {
		return nil, errors.New("event repository not available")
	}
	event, err := u.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, errors.New("event not found")
	}
	if event.PhotographerID != photographerID {
		return nil, errors.New("unauthorized")
	}

	s3KeyWatermarked := strings.Replace(input.S3Key, "/photos/", "/watermarked/", 1)

	photo := &domain.Photo{
		ID:               uuid.New(),
		EventID:          &eventID,
		S3Key:            input.S3Key,
		S3KeyWatermarked: s3KeyWatermarked,
		Filename:         input.Filename,
		SizeBytes:        input.SizeBytes,
		Width:            input.Width,
		Height:           input.Height,
	}

	if err := u.photoRepo.Create(ctx, photo); err != nil {
		return nil, err
	}
	u.userRepo.UpdateStorageUsed(ctx, photographerID, input.SizeBytes)

	return photo, nil
}

func (u *PhotoUsecase) ListEventPhotos(ctx context.Context, eventID uuid.UUID, page, limit int) ([]domain.Photo, int64, error) {
	photos, total, err := u.photoRepo.FindByEventID(ctx, eventID, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return attachPreviewURLs(photos, u.s3), total, nil
}

// ─── Delete ────────────────────────────────────────────────────────────────────

func (u *PhotoUsecase) Delete(ctx context.Context, photoID, photographerID uuid.UUID) error {
	photo, err := u.photoRepo.FindByID(ctx, photoID)
	if err != nil {
		return errors.New("photo not found")
	}

	// Ownership check: photo can be in gallery or event
	if photo.GalleryID != nil {
		gallery, err := u.galleryRepo.FindByID(ctx, *photo.GalleryID)
		if err != nil {
			return errors.New("gallery not found")
		}
		if gallery.PhotographerID != photographerID {
			return errors.New("unauthorized")
		}
	} else if photo.EventID != nil && u.eventRepo != nil {
		event, err := u.eventRepo.FindByID(ctx, *photo.EventID)
		if err != nil {
			return errors.New("event not found")
		}
		if event.PhotographerID != photographerID {
			return errors.New("unauthorized")
		}
	}

	if photo.S3Key != "" {
		u.s3.DeleteObject(photo.S3Key)
	}
	if photo.S3KeyCompressed != "" {
		u.s3.DeleteObject(photo.S3KeyCompressed)
	}
	if photo.S3KeyWatermarked != "" {
		u.s3.DeleteObject(photo.S3KeyWatermarked)
	}

	u.userRepo.UpdateStorageUsed(ctx, photographerID, -photo.SizeBytes)
	return u.photoRepo.Delete(ctx, photoID)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func attachPreviewURLs(photos []domain.Photo, s3 *infra.S3Client) []domain.Photo {
	for i := range photos {
		key := photos[i].S3KeyWatermarked
		if key == "" {
			key = photos[i].S3Key
		}
		if key != "" {
			url, err := s3.GetPresignedDownloadURL(key, 1*time.Hour)
			if err == nil {
				photos[i].PreviewURL = url
			}
		}
	}
	return photos
}
