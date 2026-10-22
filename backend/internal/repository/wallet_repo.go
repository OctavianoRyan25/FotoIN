package repository

import (
	"context"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepo struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepo {
	return &WalletRepo{db: db}
}

func (r *WalletRepo) GetBalance(ctx context.Context, photographerID uuid.UUID) (float64, error) {
	var user model.UserModel
	if err := infra.GetExecutor(ctx, r.db).Select("wallet_balance").First(&user, "id = ?", photographerID).Error; err != nil {
		return 0, err
	}
	return user.WalletBalance, nil
}

func (r *WalletRepo) AddTransaction(ctx context.Context, tx *domain.WalletTransaction) error {
	m := model.FromDomainWalletTransaction(tx)
	if err := infra.GetExecutor(ctx, r.db).Create(m).Error; err != nil {
		return err
	}
	*tx = *m.ToDomain()
	return nil
}

func (r *WalletRepo) GetTransactions(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]domain.WalletTransaction, int64, error) {
	var models []model.WalletTransactionModel
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.WalletTransactionModel{}).Where("photographer_id = ?", photographerID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return model.ToDomainWalletTransactionList(models), total, nil
}

func (r *WalletRepo) CreateWithdrawal(ctx context.Context, req *domain.WithdrawalRequest) error {
	m := model.FromDomainWithdrawalRequest(req)
	if err := infra.GetExecutor(ctx, r.db).Create(m).Error; err != nil {
		return err
	}
	*req = *m.ToDomain()
	return nil
}

func (r *WalletRepo) FindWithdrawalByID(ctx context.Context, id uuid.UUID) (*domain.WithdrawalRequest, error) {
	var m model.WithdrawalRequestModel
	if err := infra.GetExecutor(ctx, r.db).Preload("Photographer").First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *WalletRepo) ListWithdrawals(ctx context.Context, status domain.WithdrawalStatus, page, limit int) ([]domain.WithdrawalRequest, int64, error) {
	var models []model.WithdrawalRequestModel
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.WithdrawalRequestModel{}).Preload("Photographer")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return model.ToDomainWithdrawalRequestList(models), total, nil
}

func (r *WalletRepo) ListWithdrawalsByPhotographer(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]domain.WithdrawalRequest, int64, error) {
	var models []model.WithdrawalRequestModel
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&model.WithdrawalRequestModel{}).Where("photographer_id = ?", photographerID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return model.ToDomainWithdrawalRequestList(models), total, nil
}

func (r *WalletRepo) UpdateWithdrawalStatus(ctx context.Context, id uuid.UUID, status domain.WithdrawalStatus, processedBy uuid.UUID, notes string) error {
	return infra.GetExecutor(ctx, r.db).Model(&model.WithdrawalRequestModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       string(status),
			"processed_by": processedBy,
			"notes":        notes,
		}).Error
}
