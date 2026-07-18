package repository

import (
	"context"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type walletRepo struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) domain.WalletRepository {
	return &walletRepo{db: db}
}

func (r *walletRepo) GetBalance(ctx context.Context, photographerID uuid.UUID) (float64, error) {
	var user domain.User
	if err := infra.GetExecutor(ctx, r.db).Select("wallet_balance").First(&user, "id = ?", photographerID).Error; err != nil {
		return 0, err
	}
	return user.WalletBalance, nil
}

func (r *walletRepo) AddTransaction(ctx context.Context, tx *domain.WalletTransaction) error {
	return infra.GetExecutor(ctx, r.db).Create(tx).Error
}

func (r *walletRepo) GetTransactions(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]domain.WalletTransaction, int64, error) {
	var txs []domain.WalletTransaction
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&domain.WalletTransaction{}).Where("photographer_id = ?", photographerID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&txs).Error; err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

func (r *walletRepo) CreateWithdrawal(ctx context.Context, req *domain.WithdrawalRequest) error {
	return infra.GetExecutor(ctx, r.db).Create(req).Error
}

func (r *walletRepo) FindWithdrawalByID(ctx context.Context, id uuid.UUID) (*domain.WithdrawalRequest, error) {
	var req domain.WithdrawalRequest
	if err := infra.GetExecutor(ctx, r.db).Preload("Photographer").First(&req, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *walletRepo) ListWithdrawals(ctx context.Context, status domain.WithdrawalStatus, page, limit int) ([]domain.WithdrawalRequest, int64, error) {
	var reqs []domain.WithdrawalRequest
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&domain.WithdrawalRequest{}).Preload("Photographer")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&reqs).Error; err != nil {
		return nil, 0, err
	}

	return reqs, total, nil
}

func (r *walletRepo) ListWithdrawalsByPhotographer(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]domain.WithdrawalRequest, int64, error) {
	var reqs []domain.WithdrawalRequest
	var total int64

	executor := infra.GetExecutor(ctx, r.db)
	query := executor.Model(&domain.WithdrawalRequest{}).Where("photographer_id = ?", photographerID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&reqs).Error; err != nil {
		return nil, 0, err
	}

	return reqs, total, nil
}

func (r *walletRepo) UpdateWithdrawalStatus(ctx context.Context, id uuid.UUID, status domain.WithdrawalStatus, processedBy uuid.UUID, notes string) error {
	return infra.GetExecutor(ctx, r.db).Model(&domain.WithdrawalRequest{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"processed_by": processedBy,
			"notes":        notes,
		}).Error
}
