package repository

import (
	"context"
	"time"

	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(ctx context.Context, sub *domain.Subscription) error {
	m := model.FromDomainSubscription(sub)
	if err := infra.GetExecutor(ctx, r.db).Create(m).Error; err != nil {
		return err
	}
	*sub = *m.ToDomain()
	return nil
}

func (r *SubscriptionRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	var m model.SubscriptionModel
	if err := infra.GetExecutor(ctx, r.db).Preload("User").First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *SubscriptionRepo) FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*domain.Subscription, error) {
	var m model.SubscriptionModel
	err := infra.GetExecutor(ctx, r.db).
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *SubscriptionRepo) ListByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.Subscription, int64, error) {
	var models []model.SubscriptionModel
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.SubscriptionModel{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return model.ToDomainSubscriptionList(models), total, nil
}

func (r *SubscriptionRepo) ListAll(ctx context.Context, status domain.SubscriptionStatus, page, limit int) ([]domain.Subscription, int64, error) {
	var models []model.SubscriptionModel
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.SubscriptionModel{}).Preload("User")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return model.ToDomainSubscriptionList(models), total, nil
}

func (r *SubscriptionRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.SubscriptionStatus, paidAt, expiresAt *time.Time) error {
	updates := map[string]interface{}{
		"status": string(status),
	}
	if paidAt != nil {
		updates["paid_at"] = paidAt
	}
	if expiresAt != nil {
		updates["expires_at"] = expiresAt
	}
	return infra.GetExecutor(ctx, r.db).Model(&model.SubscriptionModel{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *SubscriptionRepo) UpdateQRCode(ctx context.Context, id uuid.UUID, qrURL, qrString string, expiresAt *time.Time) error {
	return infra.GetExecutor(ctx, r.db).Model(&model.SubscriptionModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"qr_code_url":    qrURL,
			"qr_code_string": qrString,
			"qr_expires_at":  expiresAt,
		}).Error
}

func (r *SubscriptionRepo) SumPaidAmount(ctx context.Context) (float64, error) {
	var sum float64
	err := infra.GetExecutor(ctx, r.db).
		Model(&model.SubscriptionModel{}).
		Where("status = ?", domain.SubscriptionStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum).Error
	return sum, err
}
