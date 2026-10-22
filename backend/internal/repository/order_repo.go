package repository

import (
	"context"
	"encoding/json"
	"time"

	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(ctx context.Context, order *domain.Order) error {
	m := model.FromDomainOrder(order)
	if err := infra.GetExecutor(ctx, r.db).Create(m).Error; err != nil {
		return err
	}
	*order = *m.ToDomain()
	return nil
}

func (r *OrderRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var m model.OrderModel
	if err := infra.GetExecutor(ctx, r.db).Preload("OrderItems.Photo").Preload("Gallery").First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *OrderRepo) FindByBuyerEmail(ctx context.Context, email string, page, limit int) ([]domain.Order, int64, error) {
	var models []model.OrderModel
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.OrderModel{}).Where("buyer_email = ?", email)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Gallery").Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return model.ToDomainOrderList(models), total, nil
}

func (r *OrderRepo) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus, gatewayID, gatewayResp string) error {
	respJSON, _ := json.Marshal(gatewayResp)
	respStr := string(respJSON)
	return infra.GetExecutor(ctx, r.db).Model(&model.OrderModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"payment_status":           string(status),
			"payment_gateway_id":       gatewayID,
			"payment_gateway_response": respStr,
		}).Error
}

func (r *OrderRepo) UpdateQRCode(ctx context.Context, id uuid.UUID, qrURL, qrString string, expiresAt *time.Time) error {
	return infra.GetExecutor(ctx, r.db).Model(&model.OrderModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"qr_code_url":    qrURL,
			"qr_code_string": qrString,
			"qr_expires_at":  expiresAt,
		}).Error
}

func (r *OrderRepo) UpdateDownloadURLs(ctx context.Context, orderID uuid.UUID, items []domain.OrderItem) error {
	executor := infra.GetExecutor(ctx, r.db)
	for _, item := range items {
		if err := executor.Model(&model.OrderItemModel{}).
			Where("id = ?", item.ID).
			Update("download_url", item.DownloadURL).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *OrderRepo) UpdateDownloadExpiresAt(ctx context.Context, orderID uuid.UUID, expiresAt time.Time) error {
	return infra.GetExecutor(ctx, r.db).Model(&model.OrderModel{}).
		Where("id = ?", orderID).
		Update("download_expires_at", expiresAt).Error
}

func (r *OrderRepo) FindByGatewayID(ctx context.Context, gatewayID string) (*domain.Order, error) {
	var m model.OrderModel
	if err := infra.GetExecutor(ctx, r.db).Preload("OrderItems").First(&m, "payment_gateway_id = ?", gatewayID).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *OrderRepo) SumPlatformFeeByGallery(ctx context.Context, galleryID uuid.UUID) (float64, error) {
	var sum float64
	err := infra.GetExecutor(ctx, r.db).Model(&model.OrderModel{}).
		Select("COALESCE(SUM(platform_fee), 0)").
		Where("gallery_id = ? AND payment_status = ?", galleryID, domain.PaymentPaid).
		Scan(&sum).Error
	return sum, err
}

func (r *OrderRepo) CountByGallery(ctx context.Context, galleryID uuid.UUID) (int64, error) {
	var count int64
	err := infra.GetExecutor(ctx, r.db).Model(&model.OrderModel{}).
		Where("gallery_id = ? AND payment_status = ?", galleryID, domain.PaymentPaid).
		Count(&count).Error
	return count, err
}
