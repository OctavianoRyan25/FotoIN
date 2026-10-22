package repository

import (
	"context"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(ctx context.Context, event *domain.Event) error {
	m := model.FromDomainEvent(event)
	if err := infra.GetExecutor(ctx, r.db).Create(m).Error; err != nil {
		return err
	}
	*event = *m.ToDomain()
	return nil
}

func (r *EventRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	var m model.EventModel
	err := infra.GetExecutor(ctx, r.db).Preload("Photographer").First(&m, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *EventRepository) FindBySlug(ctx context.Context, slug string) (*domain.Event, error) {
	var m model.EventModel
	err := infra.GetExecutor(ctx, r.db).Preload("Photographer").First(&m, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *EventRepository) Update(ctx context.Context, event *domain.Event) error {
	m := model.FromDomainEvent(event)
	if err := infra.GetExecutor(ctx, r.db).Save(m).Error; err != nil {
		return err
	}
	*event = *m.ToDomain()
	return nil
}

func (r *EventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return infra.GetExecutor(ctx, r.db).Delete(&model.EventModel{}, "id = ?", id).Error
}

func (r *EventRepository) ListByPhotographer(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]domain.Event, int64, error) {
	var models []model.EventModel
	var total int64

	offset := (page - 1) * limit
	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.EventModel{}).Where("photographer_id = ?", photographerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&models).Error
	return model.ToDomainEventList(models), total, err
}

func (r *EventRepository) ListPublic(ctx context.Context, page, limit int, search string) ([]domain.Event, int64, error) {
	var models []model.EventModel
	var total int64

	offset := (page - 1) * limit
	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.EventModel{}).Where("is_active = ?", true)

	if search != "" {
		query = query.Where("title ILIKE ? OR location ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Photographer").Order("event_date DESC").Offset(offset).Limit(limit).Find(&models).Error
	return model.ToDomainEventList(models), total, err
}
