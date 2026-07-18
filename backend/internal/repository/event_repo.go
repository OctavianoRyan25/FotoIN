package repository

import (
	"context"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) domain.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event *domain.Event) error {
	return infra.GetExecutor(ctx, r.db).Create(event).Error
}

func (r *eventRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	var event domain.Event
	err := infra.GetExecutor(ctx, r.db).Preload("Photographer").First(&event, "id = ?", id).Error
	return &event, err
}

func (r *eventRepository) FindBySlug(ctx context.Context, slug string) (*domain.Event, error) {
	var event domain.Event
	err := infra.GetExecutor(ctx, r.db).Preload("Photographer").First(&event, "slug = ?", slug).Error
	return &event, err
}

func (r *eventRepository) Update(ctx context.Context, event *domain.Event) error {
	return infra.GetExecutor(ctx, r.db).Save(event).Error
}

func (r *eventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return infra.GetExecutor(ctx, r.db).Delete(&domain.Event{}, "id = ?", id).Error
}

func (r *eventRepository) ListByPhotographer(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]domain.Event, int64, error) {
	var events []domain.Event
	var total int64

	offset := (page - 1) * limit
	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&domain.Event{}).Where("photographer_id = ?", photographerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&events).Error
	return events, total, err
}

func (r *eventRepository) ListPublic(ctx context.Context, page, limit int, search string) ([]domain.Event, int64, error) {
	var events []domain.Event
	var total int64

	offset := (page - 1) * limit
	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&domain.Event{}).Where("is_active = ?", true)

	if search != "" {
		query = query.Where("title ILIKE ? OR location ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Photographer").Order("event_date DESC").Offset(offset).Limit(limit).Find(&events).Error
	return events, total, err
}
