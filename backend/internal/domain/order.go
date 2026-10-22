package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"
	PaymentPaid    PaymentStatus = "paid"
	PaymentFailed  PaymentStatus = "failed"
	PaymentExpired PaymentStatus = "expired"
)

type Order struct {
	ID                     uuid.UUID     `json:"id"`
	BuyerEmail             string        `json:"buyer_email"`
	GalleryID              *uuid.UUID    `json:"gallery_id,omitempty"`
	EventID                *uuid.UUID    `json:"event_id,omitempty"`
	TotalAmount            float64       `json:"total_amount"`
	PlatformFee            float64       `json:"platform_fee"`
	NetAmount              float64       `json:"net_amount"`
	PaymentStatus          PaymentStatus `json:"payment_status"`
	PaymentGatewayID       *string       `json:"payment_gateway_id,omitempty"`
	PaymentGatewayResponse *string       `json:"-"`
	QRCodeURL              *string       `json:"qr_code_url,omitempty"`
	QRCodeString           *string       `json:"qr_code_string,omitempty"`
	QRExpiresAt            *time.Time    `json:"qr_expires_at,omitempty"`
	DownloadExpiresAt      *time.Time    `json:"download_expires_at,omitempty"`
	CreatedAt              time.Time     `json:"created_at"`

	Gallery    *Gallery    `json:"gallery,omitempty"`
	Event      *Event      `json:"event,omitempty"`
	OrderItems []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID          uuid.UUID `json:"id"`
	OrderID     uuid.UUID `json:"order_id"`
	PhotoID     uuid.UUID `json:"photo_id"`
	Price       float64   `json:"price"`
	DownloadURL *string   `json:"download_url,omitempty"`

	Photo *Photo `json:"photo,omitempty"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*Order, error)
	FindByBuyerEmail(ctx context.Context, email string, page, limit int) ([]Order, int64, error)
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status PaymentStatus, gatewayID, gatewayResp string) error
	UpdateQRCode(ctx context.Context, id uuid.UUID, qrURL, qrString string, expiresAt *time.Time) error
	UpdateDownloadURLs(ctx context.Context, orderID uuid.UUID, items []OrderItem) error
	UpdateDownloadExpiresAt(ctx context.Context, orderID uuid.UUID, expiresAt time.Time) error
	FindByGatewayID(ctx context.Context, gatewayID string) (*Order, error)
	SumPlatformFeeByGallery(ctx context.Context, galleryID uuid.UUID) (float64, error)
	CountByGallery(ctx context.Context, galleryID uuid.UUID) (int64, error)
}
