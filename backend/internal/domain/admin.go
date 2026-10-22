package domain

import "context"

type AdminRevenueStats struct {
	TotalRevenue        float64 `json:"total_revenue"`
	OrderCommission     float64 `json:"order_commission"`
	BasicTierRevenue    float64 `json:"basic_tier_revenue"`
	PremiumTierRevenue  float64 `json:"premium_tier_revenue"`
	SubscriptionRevenue float64 `json:"subscription_revenue"`
	TotalOrders         int64   `json:"total_orders"`
	TotalPhotographers  int64   `json:"total_photographers"`
	ActiveSubscribers   int64   `json:"active_subscribers"`
}

type AdminRepository interface {
	GetRevenueStats(ctx context.Context) (*AdminRevenueStats, error)
}
