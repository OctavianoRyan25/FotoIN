package repository

import (
	"context"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	m := model.FromDomainUser(user)
	if err := infra.GetExecutor(ctx, r.db).Create(m).Error; err != nil {
		return err
	}
	*user = *m.ToDomain()
	return nil
}

func (r *UserRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var m model.UserModel
	if err := infra.GetExecutor(ctx, r.db).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m model.UserModel
	if err := infra.GetExecutor(ctx, r.db).First(&m, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	m := model.FromDomainUser(user)
	if err := infra.GetExecutor(ctx, r.db).Save(m).Error; err != nil {
		return err
	}
	*user = *m.ToDomain()
	return nil
}

func (r *UserRepo) UpdateStorageUsed(ctx context.Context, id uuid.UUID, deltaBytes int64) error {
	return infra.GetExecutor(ctx, r.db).Model(&model.UserModel{}).
		Where("id = ?", id).
		Update("storage_used_bytes", gorm.Expr("storage_used_bytes + ?", deltaBytes)).Error
}

func (r *UserRepo) UpdateWalletBalance(ctx context.Context, id uuid.UUID, amount float64) error {
	return infra.GetExecutor(ctx, r.db).Model(&model.UserModel{}).
		Where("id = ?", id).
		Update("wallet_balance", gorm.Expr("wallet_balance + ?", amount)).Error
}

func (r *UserRepo) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	return infra.GetExecutor(ctx, r.db).Model(&model.UserModel{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).Error
}

func (r *UserRepo) CreatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	return infra.GetExecutor(ctx, r.db).Create(&model.UserModel{
		ID:           id,
		PasswordHash: passwordHash,
	}).Error
}

func (r *UserRepo) List(ctx context.Context, page, limit int, search string, role string) ([]domain.User, int64, error) {
	var models []model.UserModel
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.UserModel{})
	if role != "" && role != "all" {
		query = query.Where("role = ?", role)
	}
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return model.ToDomainUserList(models), total, nil
}

func (r *UserRepo) SuspendUser(ctx context.Context, id uuid.UUID, suspended bool) error {
	return infra.GetExecutor(ctx, r.db).Model(&model.UserModel{}).
		Where("id = ?", id).
		Update("is_suspended", suspended).Error
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return infra.GetExecutor(ctx, r.db).Where("id = ?", id).Delete(&model.UserModel{}).Error
}

func (r *UserRepo) AddOTP(ctx context.Context, otp string, expiresAt *time.Time, id uuid.UUID) error {
	return infra.GetExecutor(ctx, r.db).Model(&model.UserModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"otp_hash":       otp,
			"otp_expires_at": expiresAt,
		}).Error
}
