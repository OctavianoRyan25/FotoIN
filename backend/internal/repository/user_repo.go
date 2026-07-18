package repository

import (
	"context"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	return infra.GetExecutor(ctx, r.db).Create(user).Error
}

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := infra.GetExecutor(ctx, r.db).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := infra.GetExecutor(ctx, r.db).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user *domain.User) error {
	return infra.GetExecutor(ctx, r.db).Save(user).Error
}

func (r *userRepo) UpdateStorageUsed(ctx context.Context, id uuid.UUID, deltaBytes int64) error {
	return infra.GetExecutor(ctx, r.db).Model(&domain.User{}).
		Where("id = ?", id).
		Update("storage_used_bytes", gorm.Expr("storage_used_bytes + ?", deltaBytes)).Error
}

func (r *userRepo) UpdateWalletBalance(ctx context.Context, id uuid.UUID, amount float64) error {
	return infra.GetExecutor(ctx, r.db).Model(&domain.User{}).
		Where("id = ?", id).
		Update("wallet_balance", gorm.Expr("wallet_balance + ?", amount)).Error
}

func (r *userRepo) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	return infra.GetExecutor(ctx, r.db).Model(&domain.User{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).Error
}

func (r *userRepo) CreatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	return infra.GetExecutor(ctx, r.db).Create(&domain.User{
		ID:           id,
		PasswordHash: passwordHash,
	}).Error
}

func (r *userRepo) List(ctx context.Context, page, limit int, search string) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&domain.User{}).Where("role = ?", domain.RolePhotographer)
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepo) SuspendUser(ctx context.Context, id uuid.UUID, suspended bool) error {
	return infra.GetExecutor(ctx, r.db).Model(&domain.User{}).
		Where("id = ?", id).
		Update("is_suspended", suspended).Error
}

func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return infra.GetExecutor(ctx, r.db).Where("id = ?", id).Delete(&domain.User{}).Error
}

func (r *userRepo) AddOTP(ctx context.Context, otp string, expiresAt *time.Time, id uuid.UUID) error {
	return infra.GetExecutor(ctx, r.db).Model(&domain.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"otp_hash":       otp,
			"otp_expires_at": expiresAt,
		}).Error
}
