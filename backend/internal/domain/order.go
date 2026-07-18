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
	ID                     uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BuyerEmail             string        `gorm:"size:255;not null;index" json:"buyer_email"`
	GalleryID              *uuid.UUID    `gorm:"type:uuid;null;index" json:"gallery_id,omitempty"`
	EventID                *uuid.UUID    `gorm:"type:uuid;null;index" json:"event_id,omitempty"`
	TotalAmount            float64       `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	PlatformFee            float64       `gorm:"type:decimal(12,2);not null" json:"platform_fee"`
	NetAmount              float64       `gorm:"type:decimal(12,2);not null" json:"net_amount"`
	PaymentStatus          PaymentStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"payment_status"`
	PaymentGatewayID       *string       `gorm:"size:255;null" json:"payment_gateway_id,omitempty"`
	PaymentGatewayResponse *string       `gorm:"type:jsonb;null" json:"-"`
	// QR code fields — populated when order is created
	QRCodeURL    *string    `gorm:"size:500;null" json:"qr_code_url,omitempty"`     // hosted image (may be empty in sandbox)
	QRCodeString *string    `gorm:"type:text;null" json:"qr_code_string,omitempty"` // raw QRIS EMV string
	QRExpiresAt  *time.Time `gorm:"null" json:"qr_expires_at,omitempty"`
	// After payment is confirmed
	DownloadExpiresAt *time.Time `gorm:"null" json:"download_expires_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`

	// Relations
	Gallery    *Gallery    `gorm:"foreignKey:GalleryID" json:"gallery,omitempty"`
	Event      *Event      `gorm:"foreignKey:EventID" json:"event,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

type OrderItem struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	PhotoID     uuid.UUID `gorm:"type:uuid;not null" json:"photo_id"`
	Price       float64   `gorm:"type:decimal(12,2);not null" json:"price"`
	DownloadURL *string   `gorm:"size:500;null" json:"download_url,omitempty"` // Presigned URL after paid

	// Relations
	Photo Photo `gorm:"foreignKey:PhotoID" json:"photo,omitempty"`
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
