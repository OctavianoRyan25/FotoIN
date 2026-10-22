package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"fotoin/backend/internal/domain"

	"github.com/google/uuid"
)

type SubscriptionUsecase struct {
	subRepo        domain.SubscriptionRepository
	userRepo       domain.UserRepository
	paymentService domain.PaymentService
	txManager      domain.TxManager
}

func NewSubscriptionUsecase(
	subRepo domain.SubscriptionRepository,
	userRepo domain.UserRepository,
	paymentService domain.PaymentService,
	txManager domain.TxManager,
) *SubscriptionUsecase {
	return &SubscriptionUsecase{
		subRepo:        subRepo,
		userRepo:       userRepo,
		paymentService: paymentService,
		txManager:      txManager,
	}
}

type SubscriptionStatusResponse struct {
	User                 *domain.User         `json:"user"`
	CurrentTier          domain.SubscriptionTier `json:"current_tier"`
	ExpiresAt            *time.Time           `json:"expires_at,omitempty"`
	StorageLimitBytes    int64                `json:"storage_limit_bytes"`
	StorageUsedBytes     int64                `json:"storage_used_bytes"`
	CommissionRate       float64              `json:"commission_rate"`
	LatestSubscription   *domain.Subscription `json:"latest_subscription,omitempty"`
}

type UpgradeResponse struct {
	Subscription *domain.Subscription `json:"subscription"`
	QRCodeURL    string               `json:"qr_code_url,omitempty"`
	QRCodeString string               `json:"qr_code_string,omitempty"`
	ExpiresAt    string               `json:"expires_at"`
}

func (u *SubscriptionUsecase) GetMySubscription(ctx context.Context, userID uuid.UUID) (*SubscriptionStatusResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// If expired, revert to basic
	if user.SubscriptionTier == domain.TierPremium && user.SubscriptionExpiresAt != nil && user.SubscriptionExpiresAt.Before(time.Now()) {
		user.SubscriptionTier = domain.TierBasic
		_ = u.userRepo.Update(ctx, user)
	}

	latestSub, _ := u.subRepo.FindLatestByUserID(ctx, userID)

	return &SubscriptionStatusResponse{
		User:               user,
		CurrentTier:        user.SubscriptionTier,
		ExpiresAt:          user.SubscriptionExpiresAt,
		StorageLimitBytes:  user.StorageLimit(),
		StorageUsedBytes:   user.StorageUsedBytes,
		CommissionRate:     user.CommissionRate(),
		LatestSubscription: latestSub,
	}, nil
}

func (u *SubscriptionUsecase) UpgradeToPremium(ctx context.Context, userID uuid.UUID) (*UpgradeResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	subID := uuid.New()
	amount := domain.PremiumSubscriptionPrice

	qrPayment, err := u.paymentService.GenerateQRCode(domain.QRCodePaymentRequest{
		OrderID:       "SUB-" + subID.String()[:8],
		Amount:        amount,
		CustomerEmail: user.Email,
		Description:   "Langganan FotoIN Premium (1 Bulan)",
		PaymentMethod: domain.PaymentMethodQRIS,
	})

	var qrURL, qrString string
	var qrExpiry *time.Time
	if err == nil && qrPayment != nil {
		qrURL = qrPayment.QRCodeURL
		qrString = qrPayment.QRCodeString
		if qrPayment.ExpiresAt != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", qrPayment.ExpiresAt); err == nil {
				qrExpiry = &t
			}
		}
	} else {
		// Mock QR string for dev/testing
		qrString = fmt.Sprintf("00020101021226670016ID.CO.FOTOIN.WWW011893600002011%s52045812530336054061490005802ID5910FOTOIN6007JAKARTA6304", subID.String()[:8])
		t := time.Now().Add(30 * time.Minute)
		qrExpiry = &t
	}

	sub := &domain.Subscription{
		ID:             subID,
		UserID:         userID,
		Tier:           domain.TierPremium,
		Amount:         amount,
		DurationMonths: 1,
		Status:         domain.SubscriptionStatusPending,
		QRCodeURL:      &qrURL,
		QRCodeString:   &qrString,
		QRExpiresAt:    qrExpiry,
	}

	if err := u.subRepo.Create(ctx, sub); err != nil {
		return nil, err
	}

	expiryStr := ""
	if qrExpiry != nil {
		expiryStr = qrExpiry.Format("02 Jan 2006, 15:04 WIB")
	}

	return &UpgradeResponse{
		Subscription: sub,
		QRCodeURL:    qrURL,
		QRCodeString: qrString,
		ExpiresAt:    expiryStr,
	}, nil
}

func (u *SubscriptionUsecase) SimulatePayment(ctx context.Context, subID, userID uuid.UUID, isAdmin bool) (*domain.Subscription, error) {
	sub, err := u.subRepo.FindByID(ctx, subID)
	if err != nil {
		return nil, errors.New("subscription not found")
	}

	if !isAdmin && sub.UserID != userID {
		return nil, errors.New("unauthorized to pay for this subscription")
	}

	if sub.Status == domain.SubscriptionStatusPaid {
		return sub, nil // already paid
	}

	now := time.Now()
	var newExpiresAt time.Time

	user, err := u.userRepo.FindByID(ctx, sub.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.SubscriptionTier == domain.TierPremium && user.SubscriptionExpiresAt != nil && user.SubscriptionExpiresAt.After(now) {
		newExpiresAt = user.SubscriptionExpiresAt.AddDate(0, sub.DurationMonths, 0)
	} else {
		newExpiresAt = now.AddDate(0, sub.DurationMonths, 0)
	}

	err = u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := u.subRepo.UpdateStatus(txCtx, sub.ID, domain.SubscriptionStatusPaid, &now, &newExpiresAt); err != nil {
			return err
		}

		user.SubscriptionTier = domain.TierPremium
		user.SubscriptionExpiresAt = &newExpiresAt
		return u.userRepo.Update(txCtx, user)
	})

	if err != nil {
		return nil, err
	}

	sub.Status = domain.SubscriptionStatusPaid
	sub.PaidAt = &now
	sub.ExpiresAt = &newExpiresAt
	return sub, nil
}

func (u *SubscriptionUsecase) ListAllSubscriptions(ctx context.Context, status domain.SubscriptionStatus, page, limit int) ([]domain.Subscription, int64, error) {
	return u.subRepo.ListAll(ctx, status, page, limit)
}

func (u *SubscriptionUsecase) AdminGrantPremium(ctx context.Context, targetUserID uuid.UUID, months int) (*domain.User, error) {
	if months <= 0 {
		months = 1
	}

	user, err := u.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return nil, errors.New("target user not found")
	}

	now := time.Now()
	var newExpiresAt time.Time
	if user.SubscriptionTier == domain.TierPremium && user.SubscriptionExpiresAt != nil && user.SubscriptionExpiresAt.After(now) {
		newExpiresAt = user.SubscriptionExpiresAt.AddDate(0, months, 0)
	} else {
		newExpiresAt = now.AddDate(0, months, 0)
	}

	err = u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		user.SubscriptionTier = domain.TierPremium
		user.SubscriptionExpiresAt = &newExpiresAt
		if err := u.userRepo.Update(txCtx, user); err != nil {
			return err
		}

		sub := &domain.Subscription{
			ID:             uuid.New(),
			UserID:         targetUserID,
			Tier:           domain.TierPremium,
			Amount:         0,
			DurationMonths: months,
			Status:         domain.SubscriptionStatusPaid,
			PaidAt:         &now,
			ExpiresAt:      &newExpiresAt,
		}
		return u.subRepo.Create(txCtx, sub)
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}
