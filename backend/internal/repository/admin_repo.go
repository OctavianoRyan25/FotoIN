package repository

import (
	"context"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/model"

	"gorm.io/gorm"
)

type AdminRepo struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepo {
	return &AdminRepo{db: db}
}

func (r *AdminRepo) GetRevenueStats(ctx context.Context) (*domain.AdminRevenueStats, error) {
	executor := infra.GetExecutor(ctx, r.db)
	var stats domain.AdminRevenueStats

	if err := executor.Model(&model.OrderModel{}).
		Select("COALESCE(SUM(platform_fee), 0) as order_commission, COUNT(*) as total_orders").
		Where("payment_status = ?", domain.PaymentPaid).
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	var basicGal, premiumGal float64
	if err := executor.Raw(`
		SELECT COALESCE(SUM(o.platform_fee), 0)
		FROM orders o
		JOIN galleries g ON o.gallery_id = g.id
		JOIN users u ON g.photographer_id = u.id
		WHERE o.payment_status = ? AND u.subscription_tier = ?
	`, domain.PaymentPaid, domain.TierBasic).Scan(&basicGal).Error; err != nil {
		return nil, err
	}

	if err := executor.Raw(`
		SELECT COALESCE(SUM(o.platform_fee), 0)
		FROM orders o
		JOIN galleries g ON o.gallery_id = g.id
		JOIN users u ON g.photographer_id = u.id
		WHERE o.payment_status = ? AND u.subscription_tier = ?
	`, domain.PaymentPaid, domain.TierPremium).Scan(&premiumGal).Error; err != nil {
		return nil, err
	}

	var basicEv, premiumEv float64
	if err := executor.Raw(`
		SELECT COALESCE(SUM(o.platform_fee), 0)
		FROM orders o
		JOIN events e ON o.event_id = e.id
		JOIN users u ON e.photographer_id = u.id
		WHERE o.payment_status = ? AND u.subscription_tier = ?
	`, domain.PaymentPaid, domain.TierBasic).Scan(&basicEv).Error; err != nil {
		return nil, err
	}

	if err := executor.Raw(`
		SELECT COALESCE(SUM(o.platform_fee), 0)
		FROM orders o
		JOIN events e ON o.event_id = e.id
		JOIN users u ON e.photographer_id = u.id
		WHERE o.payment_status = ? AND u.subscription_tier = ?
	`, domain.PaymentPaid, domain.TierPremium).Scan(&premiumEv).Error; err != nil {
		return nil, err
	}

	stats.BasicTierRevenue = basicGal + basicEv
	stats.PremiumTierRevenue = premiumGal + premiumEv

	if err := executor.Model(&model.SubscriptionModel{}).
		Where("status = ?", domain.SubscriptionStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.SubscriptionRevenue).Error; err != nil {
		return nil, err
	}

	stats.TotalRevenue = stats.OrderCommission + stats.SubscriptionRevenue

	if err := executor.Model(&model.UserModel{}).
		Where("role = ?", domain.RolePhotographer).
		Count(&stats.TotalPhotographers).Error; err != nil {
		return nil, err
	}

	if err := executor.Model(&model.UserModel{}).
		Where("subscription_tier = ? AND (subscription_expires_at IS NULL OR subscription_expires_at > NOW())", domain.TierPremium).
		Count(&stats.ActiveSubscribers).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
