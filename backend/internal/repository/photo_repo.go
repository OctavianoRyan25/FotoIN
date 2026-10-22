package repository

import (
	"context"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PhotoRepo struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) *PhotoRepo {
	return &PhotoRepo{db: db}
}

func (r *PhotoRepo) Create(ctx context.Context, photo *domain.Photo) error {
	m := model.FromDomainPhoto(photo)
	if err := infra.GetExecutor(ctx, r.db).Create(m).Error; err != nil {
		return err
	}
	*photo = *m.ToDomain()
	return nil
}

func (r *PhotoRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Photo, error) {
	var m model.PhotoModel
	if err := infra.GetExecutor(ctx, r.db).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *PhotoRepo) FindByGalleryID(ctx context.Context, galleryID uuid.UUID, page, limit int) ([]domain.Photo, int64, error) {
	var models []model.PhotoModel
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.PhotoModel{}).Where("gallery_id = ?", galleryID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return model.ToDomainPhotoList(models), total, nil
}

func (r *PhotoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return infra.GetExecutor(ctx, r.db).Delete(&model.PhotoModel{}, "id = ?", id).Error
}

func (r *PhotoRepo) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Photo, error) {
	var models []model.PhotoModel
	if err := infra.GetExecutor(ctx, r.db).Where("id IN ?", ids).Find(&models).Error; err != nil {
		return nil, err
	}
	return model.ToDomainPhotoList(models), nil
}

func (r *PhotoRepo) BulkDelete(ctx context.Context, galleryID uuid.UUID) error {
	return infra.GetExecutor(ctx, r.db).Where("gallery_id = ?", galleryID).Delete(&model.PhotoModel{}).Error
}

func (r *PhotoRepo) FindByEventID(ctx context.Context, eventID uuid.UUID, page, limit int) ([]domain.Photo, int64, error) {
	var models []model.PhotoModel
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.PhotoModel{}).Where("event_id = ?", eventID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return model.ToDomainPhotoList(models), total, nil
}
