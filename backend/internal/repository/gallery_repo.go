package repository

import (
	"context"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type galleryRepo struct {
	db *gorm.DB
}

func NewGalleryRepository(db *gorm.DB) domain.GalleryRepository {
	return &galleryRepo{db: db}
}

func (r *galleryRepo) Create(ctx context.Context, gallery *domain.Gallery) error {
	m := model.FromDomainGallery(gallery)
	if err := infra.GetExecutor(ctx, r.db).Create(m).Error; err != nil {
		return err
	}
	*gallery = *m.ToDomain()
	return nil
}

func (r *galleryRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Gallery, error) {
	var m model.GalleryModel
	executor := infra.GetExecutor(ctx, r.db)
	if err := executor.Preload("Photographer").First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	gallery := m.ToDomain()

	var count int64
	executor.Model(&model.PhotoModel{}).Where("gallery_id = ?", id).Count(&count)
	gallery.PhotoCount = int(count)

	return gallery, nil
}

func (r *galleryRepo) FindBySlug(ctx context.Context, slug string) (*domain.Gallery, error) {
	var m model.GalleryModel
	executor := infra.GetExecutor(ctx, r.db)
	if err := executor.Preload("Photographer").First(&m, "slug = ? AND is_active = true AND (expires_at IS NULL OR expires_at > NOW())", slug).Error; err != nil {
		return nil, err
	}
	gallery := m.ToDomain()

	var count int64
	executor.Model(&model.PhotoModel{}).Where("gallery_id = ?", gallery.ID).Count(&count)
	gallery.PhotoCount = int(count)

	return gallery, nil
}

func (r *galleryRepo) FindByPhotographerID(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]domain.Gallery, int64, error) {
	var models []model.GalleryModel
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.GalleryModel{}).Where("photographer_id = ?", photographerID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	galleries := model.ToDomainGalleryList(models)
	for i := range galleries {
		var count int64
		executor.Model(&model.PhotoModel{}).Where("gallery_id = ?", galleries[i].ID).Count(&count)
		galleries[i].PhotoCount = int(count)
	}

	return galleries, total, nil
}

func (r *galleryRepo) Update(ctx context.Context, gallery *domain.Gallery) error {
	m := model.FromDomainGallery(gallery)
	if err := infra.GetExecutor(ctx, r.db).Save(m).Error; err != nil {
		return err
	}
	*gallery = *m.ToDomain()
	return nil
}

func (r *galleryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return infra.GetExecutor(ctx, r.db).Delete(&model.GalleryModel{}, "id = ?", id).Error
}

func (r *galleryRepo) ListPublic(ctx context.Context, page, limit int) ([]domain.Gallery, int64, error) {
	var models []model.GalleryModel
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.GalleryModel{}).Where("is_active = true AND (expires_at IS NULL OR expires_at > NOW())").Preload("Photographer")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return model.ToDomainGalleryList(models), total, nil
}

func (r *galleryRepo) CountPhotos(ctx context.Context, galleryID uuid.UUID) (int, error) {
	var count int64
	err := infra.GetExecutor(ctx, r.db).Model(&model.PhotoModel{}).Where("gallery_id = ?", galleryID).Count(&count).Error
	return int(count), err
}
