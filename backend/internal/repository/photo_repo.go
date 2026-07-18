package repository

import (
	"context"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type photoRepo struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) domain.PhotoRepository {
	return &photoRepo{db: db}
}

func (r *photoRepo) Create(ctx context.Context, photo *domain.Photo) error {
	return infra.GetExecutor(ctx, r.db).Create(photo).Error
}

func (r *photoRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Photo, error) {
	var photo domain.Photo
	if err := infra.GetExecutor(ctx, r.db).First(&photo, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &photo, nil
}

func (r *photoRepo) FindByGalleryID(ctx context.Context, galleryID uuid.UUID, page, limit int) ([]domain.Photo, int64, error) {
	var photos []domain.Photo
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&domain.Photo{}).Where("gallery_id = ?", galleryID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&photos).Error; err != nil {
		return nil, 0, err
	}

	return photos, total, nil
}

func (r *photoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return infra.GetExecutor(ctx, r.db).Delete(&domain.Photo{}, "id = ?", id).Error
}

func (r *photoRepo) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Photo, error) {
	var photos []domain.Photo
	if err := infra.GetExecutor(ctx, r.db).Where("id IN ?", ids).Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (r *photoRepo) BulkDelete(ctx context.Context, galleryID uuid.UUID) error {
	return infra.GetExecutor(ctx, r.db).Where("gallery_id = ?", galleryID).Delete(&domain.Photo{}).Error
}

func (r *photoRepo) FindByEventID(ctx context.Context, eventID uuid.UUID, page, limit int) ([]domain.Photo, int64, error) {
	var photos []domain.Photo
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&domain.Photo{}).Where("event_id = ?", eventID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&photos).Error; err != nil {
		return nil, 0, err
	}

	return photos, total, nil
}

