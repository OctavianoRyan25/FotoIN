package model

import (
	"time"

	"fotoin/backend/internal/domain"

	"github.com/google/uuid"
)

type OrderModel struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BuyerEmail             string     `gorm:"size:255;not null;index"`
	GalleryID              *uuid.UUID `gorm:"type:uuid;null;index"`
	EventID                *uuid.UUID `gorm:"type:uuid;null;index"`
	TotalAmount            float64    `gorm:"type:decimal(12,2);not null"`
	PlatformFee            float64    `gorm:"type:decimal(12,2);not null"`
	NetAmount              float64    `gorm:"type:decimal(12,2);not null"`
	PaymentStatus          string     `gorm:"type:varchar(20);not null;default:'pending'"`
	PaymentGatewayID       *string    `gorm:"size:255;null"`
	PaymentGatewayResponse *string    `gorm:"type:jsonb;null"`
	QRCodeURL              *string    `gorm:"size:500;null"`
	QRCodeString           *string    `gorm:"type:text;null"`
	QRExpiresAt            *time.Time `gorm:"null"`
	DownloadExpiresAt      *time.Time `gorm:"null"`
	CreatedAt              time.Time

	// Relations
	Gallery    *GalleryModel    `gorm:"foreignKey:GalleryID"`
	Event      *EventModel      `gorm:"foreignKey:EventID"`
	OrderItems []OrderItemModel `gorm:"foreignKey:OrderID"`
}

func (OrderModel) TableName() string {
	return "orders"
}

type OrderItemModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null;index"`
	PhotoID     uuid.UUID `gorm:"type:uuid;not null"`
	Price       float64   `gorm:"type:decimal(12,2);not null"`
	DownloadURL *string   `gorm:"size:500;null"`

	// Relations
	Photo PhotoModel `gorm:"foreignKey:PhotoID"`
}

func (OrderItemModel) TableName() string {
	return "order_items"
}

func (m *OrderItemModel) ToDomain() *domain.OrderItem {
	if m == nil {
		return nil
	}
	item := &domain.OrderItem{
		ID:          m.ID,
		OrderID:     m.OrderID,
		PhotoID:     m.PhotoID,
		Price:       m.Price,
		DownloadURL: m.DownloadURL,
	}
	if m.Photo.ID != uuid.Nil {
		item.Photo = m.Photo.ToDomain()
	}
	return item
}

func FromDomainOrderItem(item *domain.OrderItem) *OrderItemModel {
	if item == nil {
		return nil
	}
	return &OrderItemModel{
		ID:          item.ID,
		OrderID:     item.OrderID,
		PhotoID:     item.PhotoID,
		Price:       item.Price,
		DownloadURL: item.DownloadURL,
	}
}

func ToDomainOrderItemList(models []OrderItemModel) []domain.OrderItem {
	res := make([]domain.OrderItem, len(models))
	for i, m := range models {
		res[i] = *m.ToDomain()
	}
	return res
}

func (m *OrderModel) ToDomain() *domain.Order {
	if m == nil {
		return nil
	}
	o := &domain.Order{
		ID:                     m.ID,
		BuyerEmail:             m.BuyerEmail,
		GalleryID:              m.GalleryID,
		EventID:                m.EventID,
		TotalAmount:            m.TotalAmount,
		PlatformFee:            m.PlatformFee,
		NetAmount:              m.NetAmount,
		PaymentStatus:          domain.PaymentStatus(m.PaymentStatus),
		PaymentGatewayID:       m.PaymentGatewayID,
		PaymentGatewayResponse: m.PaymentGatewayResponse,
		QRCodeURL:              m.QRCodeURL,
		QRCodeString:           m.QRCodeString,
		QRExpiresAt:            m.QRExpiresAt,
		DownloadExpiresAt:      m.DownloadExpiresAt,
		CreatedAt:              m.CreatedAt,
		OrderItems:             ToDomainOrderItemList(m.OrderItems),
	}

	if m.Gallery != nil && m.Gallery.ID != uuid.Nil {
		o.Gallery = m.Gallery.ToDomain()
	}
	if m.Event != nil && m.Event.ID != uuid.Nil {
		o.Event = m.Event.ToDomain()
	}

	return o
}

func FromDomainOrder(o *domain.Order) *OrderModel {
	if o == nil {
		return nil
	}
	items := make([]OrderItemModel, len(o.OrderItems))
	for i, item := range o.OrderItems {
		items[i] = *FromDomainOrderItem(&item)
	}

	return &OrderModel{
		ID:                     o.ID,
		BuyerEmail:             o.BuyerEmail,
		GalleryID:              o.GalleryID,
		EventID:                o.EventID,
		TotalAmount:            o.TotalAmount,
		PlatformFee:            o.PlatformFee,
		NetAmount:              o.NetAmount,
		PaymentStatus:          string(o.PaymentStatus),
		PaymentGatewayID:       o.PaymentGatewayID,
		PaymentGatewayResponse: o.PaymentGatewayResponse,
		QRCodeURL:              o.QRCodeURL,
		QRCodeString:           o.QRCodeString,
		QRExpiresAt:            o.QRExpiresAt,
		DownloadExpiresAt:      o.DownloadExpiresAt,
		CreatedAt:              o.CreatedAt,
		OrderItems:             items,
	}
}

func ToDomainOrderList(models []OrderModel) []domain.Order {
	res := make([]domain.Order, len(models))
	for i, m := range models {
		res[i] = *m.ToDomain()
	}
	return res
}
