package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type SubscriptionTier string
type UserRole string

const (
	TierBasic   SubscriptionTier = "basic"
	TierPremium SubscriptionTier = "premium"

	RolePhotographer UserRole = "photographer"
	RoleSuperAdmin   UserRole = "super_admin"
	RoleClient       UserRole = "client"
)

// Business constants
const (
	BasicStorageLimitBytes   = int64(3 * 1024 * 1024 * 1024)   // 3 GB
	PremiumStorageLimitBytes = int64(100 * 1024 * 1024 * 1024) // 100 GB
	BasicRetentionDays       = 14
	MinPricePerPhoto         = float64(10000) // IDR 10.000
	BasicCommissionRate      = float64(0.20)  // 20%
	PremiumCommissionRate    = float64(0.04)  // 4%
)

type User struct {
	ID                    uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name                  string           `gorm:"size:255;not null" json:"name"`
	Email                 string           `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash          string           `gorm:"size:255;not null" json:"-"`
	OTPHash               string           `gorm:"size:255;null" json:"-"`
	OTPExpiresAt          *time.Time       `gorm:"null" json:"-"`
	Role                  UserRole         `gorm:"type:varchar(20);not null;default:'photographer'" json:"role"`
	SubscriptionTier      SubscriptionTier `gorm:"type:varchar(20);not null;default:'basic'" json:"subscription_tier"`
	SubscriptionExpiresAt *time.Time       `gorm:"null" json:"subscription_expires_at,omitempty"`
	IsVerified            bool             `gorm:"default:false" json:"is_verified"`
	IsSuspended           bool             `gorm:"default:false" json:"is_suspended"`
	StorageUsedBytes      int64            `gorm:"default:0" json:"storage_used_bytes"`
	WalletBalance         float64          `gorm:"type:decimal(14,2);default:0" json:"wallet_balance"`
	CreatedAt             time.Time        `json:"created_at"`
	UpdatedAt             time.Time        `json:"updated_at"`

	HasPassword bool `gorm:"-" json:"has_password"`

	// Relations
	Galleries []Gallery `gorm:"foreignKey:PhotographerID" json:"-"`
}

func (u *User) CommissionRate() float64 {
	if u.SubscriptionTier == TierPremium {
		return PremiumCommissionRate
	}
	return BasicCommissionRate
}

func (u *User) StorageLimit() int64 {
	if u.SubscriptionTier == TierPremium {
		return PremiumStorageLimitBytes
	}
	return BasicStorageLimitBytes
}

func (u *User) StorageUsedGB() float64 {
	return float64(u.StorageUsedBytes) / (1024 * 1024 * 1024)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	UpdateStorageUsed(ctx context.Context, id uuid.UUID, deltaBytes int64) error
	UpdateWalletBalance(ctx context.Context, id uuid.UUID, amount float64) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	CreatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
	AddOTP(ctx context.Context, otp string, expiresAt *time.Time, id uuid.UUID) error
	List(ctx context.Context, page, limit int, search string) ([]User, int64, error)
	SuspendUser(ctx context.Context, id uuid.UUID, suspended bool) error
	Delete(ctx context.Context, id uuid.UUID) error
}
