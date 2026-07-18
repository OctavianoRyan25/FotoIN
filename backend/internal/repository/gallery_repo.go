package repository

import (
	"context"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"

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
	return infra.GetExecutor(ctx, r.db).Create(gallery).Error
}

func (r *galleryRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Gallery, error) {
	var gallery domain.Gallery
	executor := infra.GetExecutor(ctx, r.db)
	if err := executor.Preload("Photographer").First(&gallery, "id = ?", id).Error; err != nil {
		return nil, err
	}
	// Count photos
	var count int64
	executor.Model(&domain.Photo{}).Where("gallery_id = ?", id).Count(&count)
	gallery.PhotoCount = int(count)

	return &gallery, nil
}

func (r *galleryRepo) FindBySlug(ctx context.Context, slug string) (*domain.Gallery, error) {
	var gallery domain.Gallery
	executor := infra.GetExecutor(ctx, r.db)
	if err := executor.Preload("Photographer").First(&gallery, "slug = ? AND is_active = true", slug).Error; err != nil {
		return nil, err
	}
	var count int64
	executor.Model(&domain.Photo{}).Where("gallery_id = ?", gallery.ID).Count(&count)
	gallery.PhotoCount = int(count)

	return &gallery, nil
}

func (r *galleryRepo) FindByPhotographerID(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]domain.Gallery, int64, error) {
	var galleries []domain.Gallery
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&domain.Gallery{}).Where("photographer_id = ?", photographerID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&galleries).Error; err != nil {
		return nil, 0, err
	}

	// Enrich with photo counts
	for i := range galleries {
		var count int64
		executor.Model(&domain.Photo{}).Where("gallery_id = ?", galleries[i].ID).Count(&count)
		galleries[i].PhotoCount = int(count)
	}

	return galleries, total, nil
}

func (r *galleryRepo) Update(ctx context.Context, gallery *domain.Gallery) error {
	return infra.GetExecutor(ctx, r.db).Save(gallery).Error
}

func (r *galleryRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return infra.GetExecutor(ctx, r.db).Delete(&domain.Gallery{}, "id = ?", id).Error
}

func (r *galleryRepo) ListPublic(ctx context.Context, page, limit int) ([]domain.Gallery, int64, error) {
	var galleries []domain.Gallery
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&domain.Gallery{}).Where("is_active = true").Preload("Photographer")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&galleries).Error; err != nil {
		return nil, 0, err
	}

	return galleries, total, nil
}

func (r *galleryRepo) CountPhotos(ctx context.Context, galleryID uuid.UUID) (int, error) {
	var count int64
	err := infra.GetExecutor(ctx, r.db).Model(&domain.Photo{}).Where("gallery_id = ?", galleryID).Count(&count).Error
	return int(count), err
}
