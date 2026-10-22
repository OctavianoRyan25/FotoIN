package model

import (
	"time"

	"fotoin/backend/internal/domain"

	"github.com/google/uuid"
)

type SubscriptionModel struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null;index"`
	Tier             string     `gorm:"type:varchar(20);not null"`
	Amount           float64    `gorm:"type:decimal(12,2);not null"`
	DurationMonths   int        `gorm:"not null;default:1"`
	Status           string     `gorm:"type:varchar(20);not null;default:'pending'"`
	PaymentGatewayID *string    `gorm:"size:255;null"`
	QRCodeURL        *string    `gorm:"size:500;null"`
	QRCodeString     *string    `gorm:"type:text;null"`
	QRExpiresAt      *time.Time `gorm:"null"`
	PaidAt           *time.Time `gorm:"null"`
	ExpiresAt        *time.Time `gorm:"null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time

	// Relations
	User UserModel `gorm:"foreignKey:UserID"`
}

func (SubscriptionModel) TableName() string {
	return "subscriptions"
}

func (m *SubscriptionModel) ToDomain() *domain.Subscription {
	if m == nil {
		return nil
	}
	sub := &domain.Subscription{
		ID:               m.ID,
		UserID:           m.UserID,
		Tier:             domain.SubscriptionTier(m.Tier),
		Amount:           m.Amount,
		DurationMonths:   m.DurationMonths,
		Status:           domain.SubscriptionStatus(m.Status),
		PaymentGatewayID: m.PaymentGatewayID,
		QRCodeURL:        m.QRCodeURL,
		QRCodeString:     m.QRCodeString,
		QRExpiresAt:      m.QRExpiresAt,
		PaidAt:           m.PaidAt,
		ExpiresAt:        m.ExpiresAt,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}

	if m.User.ID != uuid.Nil {
		sub.User = m.User.ToDomain()
	}

	return sub
}

func FromDomainSubscription(s *domain.Subscription) *SubscriptionModel {
	if s == nil {
		return nil
	}
	return &SubscriptionModel{
		ID:               s.ID,
		UserID:           s.UserID,
		Tier:             string(s.Tier),
		Amount:           s.Amount,
		DurationMonths:   s.DurationMonths,
		Status:           string(s.Status),
		PaymentGatewayID: s.PaymentGatewayID,
		QRCodeURL:        s.QRCodeURL,
		QRCodeString:     s.QRCodeString,
		QRExpiresAt:      s.QRExpiresAt,
		PaidAt:           s.PaidAt,
		ExpiresAt:        s.ExpiresAt,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}

func ToDomainSubscriptionList(models []SubscriptionModel) []domain.Subscription {
	res := make([]domain.Subscription, len(models))
	for i, m := range models {
		res[i] = *m.ToDomain()
	}
	return res
}
