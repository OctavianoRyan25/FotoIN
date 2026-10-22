package model

import (
	"time"

	"fotoin/backend/internal/domain"

	"github.com/google/uuid"
)

type UserModel struct {
	ID                    uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name                  string     `gorm:"size:255;not null"`
	Email                 string     `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash          string     `gorm:"size:255;not null"`
	OTPHash               string     `gorm:"size:255;null"`
	OTPExpiresAt          *time.Time `gorm:"null"`
	Role                  string     `gorm:"type:varchar(20);not null;default:'photographer'"`
	SubscriptionTier      string     `gorm:"type:varchar(20);not null;default:'basic'"`
	SubscriptionExpiresAt *time.Time `gorm:"null"`
	IsVerified            bool       `gorm:"default:false"`
	IsSuspended           bool       `gorm:"default:false"`
	StorageUsedBytes      int64      `gorm:"default:0"`
	WalletBalance         float64    `gorm:"type:decimal(14,2);default:0"`
	CreatedAt             time.Time
	UpdatedAt             time.Time

	// Relations
	Galleries []GalleryModel `gorm:"foreignKey:PhotographerID"`
}

func (UserModel) TableName() string {
	return "users"
}

func (m *UserModel) ToDomain() *domain.User {
	if m == nil {
		return nil
	}
	return &domain.User{
		ID:                    m.ID,
		Name:                  m.Name,
		Email:                 m.Email,
		PasswordHash:          m.PasswordHash,
		OTPHash:               m.OTPHash,
		OTPExpiresAt:          m.OTPExpiresAt,
		Role:                  domain.UserRole(m.Role),
		SubscriptionTier:      domain.SubscriptionTier(m.SubscriptionTier),
		SubscriptionExpiresAt: m.SubscriptionExpiresAt,
		IsVerified:            m.IsVerified,
		IsSuspended:           m.IsSuspended,
		StorageUsedBytes:      m.StorageUsedBytes,
		WalletBalance:         m.WalletBalance,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
		HasPassword:           m.PasswordHash != "",
	}
}

func FromDomainUser(u *domain.User) *UserModel {
	if u == nil {
		return nil
	}
	return &UserModel{
		ID:                    u.ID,
		Name:                  u.Name,
		Email:                 u.Email,
		PasswordHash:          u.PasswordHash,
		OTPHash:               u.OTPHash,
		OTPExpiresAt:          u.OTPExpiresAt,
		Role:                  string(u.Role),
		SubscriptionTier:      string(u.SubscriptionTier),
		SubscriptionExpiresAt: u.SubscriptionExpiresAt,
		IsVerified:            u.IsVerified,
		IsSuspended:           u.IsSuspended,
		StorageUsedBytes:      u.StorageUsedBytes,
		WalletBalance:         u.WalletBalance,
		CreatedAt:             u.CreatedAt,
		UpdatedAt:             u.UpdatedAt,
	}
}

func ToDomainUserList(models []UserModel) []domain.User {
	res := make([]domain.User, len(models))
	for i, m := range models {
		res[i] = *m.ToDomain()
	}
	return res
}
