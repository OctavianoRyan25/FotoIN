package repository

import (
	"context"
	"encoding/json"
	"time"

	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(ctx context.Context, order *domain.Order) error {
	return infra.GetExecutor(ctx, r.db).Create(order).Error
}

func (r *orderRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	if err := infra.GetExecutor(ctx, r.db).Preload("OrderItems.Photo").Preload("Gallery").First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepo) FindByBuyerEmail(ctx context.Context, email string, page, limit int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64
	
	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&domain.Order{}).Where("buyer_email = ?", email)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Gallery").Offset(offset).Limit(limit).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepo) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus, gatewayID, gatewayResp string) error {
	respJSON, _ := json.Marshal(gatewayResp)
	respStr := string(respJSON)
	return infra.GetExecutor(ctx, r.db).Model(&domain.Order{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"payment_status":           status,
			"payment_gateway_id":       gatewayID,
			"payment_gateway_response": respStr,
		}).Error
}

func (r *orderRepo) UpdateQRCode(ctx context.Context, id uuid.UUID, qrURL, qrString string, expiresAt *time.Time) error {
	return infra.GetExecutor(ctx, r.db).Model(&domain.Order{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"qr_code_url":    qrURL,
			"qr_code_string": qrString,
			"qr_expires_at":  expiresAt,
		}).Error
}

func (r *orderRepo) UpdateDownloadURLs(ctx context.Context, orderID uuid.UUID, items []domain.OrderItem) error {
	executor := infra.GetExecutor(ctx, r.db)
	for _, item := range items {
		if err := executor.Model(&domain.OrderItem{}).
			Where("id = ?", item.ID).
			Update("download_url", item.DownloadURL).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *orderRepo) UpdateDownloadExpiresAt(ctx context.Context, orderID uuid.UUID, expiresAt time.Time) error {
	return infra.GetExecutor(ctx, r.db).Model(&domain.Order{}).
		Where("id = ?", orderID).
		Update("download_expires_at", expiresAt).Error
}

func (r *orderRepo) FindByGatewayID(ctx context.Context, gatewayID string) (*domain.Order, error) {
	var order domain.Order
	if err := infra.GetExecutor(ctx, r.db).Preload("OrderItems").First(&order, "payment_gateway_id = ?", gatewayID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepo) SumPlatformFeeByGallery(ctx context.Context, galleryID uuid.UUID) (float64, error) {
	var sum float64
	err := infra.GetExecutor(ctx, r.db).Model(&domain.Order{}).
		Select("COALESCE(SUM(platform_fee), 0)").
		Where("gallery_id = ? AND payment_status = ?", galleryID, domain.PaymentPaid).
		Scan(&sum).Error
	return sum, err
}

func (r *orderRepo) CountByGallery(ctx context.Context, galleryID uuid.UUID) (int64, error) {
	var count int64
	err := infra.GetExecutor(ctx, r.db).Model(&domain.Order{}).
		Where("gallery_id = ? AND payment_status = ?", galleryID, domain.PaymentPaid).
		Count(&count).Error
	return count, err
}
