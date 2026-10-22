package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionStatusPending SubscriptionStatus = "pending"
	SubscriptionStatusPaid    SubscriptionStatus = "paid"
	SubscriptionStatusExpired SubscriptionStatus = "expired"
)

const PremiumSubscriptionPrice = float64(149000) // Rp 149.000 / month

type Subscription struct {
	ID               uuid.UUID          `json:"id"`
	UserID           uuid.UUID          `json:"user_id"`
	Tier             SubscriptionTier   `json:"tier"`
	Amount           float64            `json:"amount"`
	DurationMonths   int                `json:"duration_months"`
	Status           SubscriptionStatus `json:"status"`
	PaymentGatewayID *string            `json:"payment_gateway_id,omitempty"`
	QRCodeURL        *string            `json:"qr_code_url,omitempty"`
	QRCodeString     *string            `json:"qr_code_string,omitempty"`
	QRExpiresAt      *time.Time         `json:"qr_expires_at,omitempty"`
	PaidAt           *time.Time         `json:"paid_at,omitempty"`
	ExpiresAt        *time.Time         `json:"expires_at,omitempty"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`

	User *User `json:"user,omitempty"`
}

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *Subscription) error
	FindByID(ctx context.Context, id uuid.UUID) (*Subscription, error)
	FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*Subscription, error)
	ListByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]Subscription, int64, error)
	ListAll(ctx context.Context, status SubscriptionStatus, page, limit int) ([]Subscription, int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status SubscriptionStatus, paidAt, expiresAt *time.Time) error
	UpdateQRCode(ctx context.Context, id uuid.UUID, qrURL, qrString string, expiresAt *time.Time) error
	SumPaidAmount(ctx context.Context) (float64, error)
}
