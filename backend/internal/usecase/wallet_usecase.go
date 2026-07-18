package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"fotoin/backend/internal/domain"

	"github.com/google/uuid"
)

type WalletUsecase struct {
	walletRepo domain.WalletRepository
	userRepo   domain.UserRepository
	txManager  domain.TxManager
}

func NewWalletUsecase(walletRepo domain.WalletRepository, userRepo domain.UserRepository, txManager domain.TxManager) *WalletUsecase {
	return &WalletUsecase{walletRepo: walletRepo, userRepo: userRepo, txManager: txManager}
}

type WithdrawalInput struct {
	Amount        float64 `json:"amount" validate:"required,min=50000"` // Min withdraw Rp 50.000
	BankName      string  `json:"bank_name" validate:"required"`
	AccountNumber string  `json:"account_number" validate:"required"`
	AccountHolder string  `json:"account_holder" validate:"required"`
}

type WalletSummary struct {
	Balance         float64                    `json:"balance"`
	Transactions    []domain.WalletTransaction `json:"transactions"`
	TotalPages      int                        `json:"total_pages"`
	Total           int64                      `json:"total"`
}

func (u *WalletUsecase) GetSummary(ctx context.Context, photographerID uuid.UUID, page, limit int) (*WalletSummary, error) {
	balance, err := u.walletRepo.GetBalance(ctx, photographerID)
	if err != nil {
		return nil, err
	}

	txs, total, err := u.walletRepo.GetTransactions(ctx, photographerID, page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &WalletSummary{
		Balance:      balance,
		Transactions: txs,
		TotalPages:   totalPages,
		Total:        total,
	}, nil
}

func (u *WalletUsecase) RequestWithdrawal(ctx context.Context, photographerID uuid.UUID, input WithdrawalInput) (*domain.WithdrawalRequest, error) {
	user, err := u.userRepo.FindByID(ctx, photographerID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.WalletBalance < input.Amount {
		return nil, fmt.Errorf("insufficient balance (current: Rp %.0f, requested: Rp %.0f)",
			user.WalletBalance, input.Amount)
	}

	// Check for existing pending withdrawal
	existingReqs, _, _ := u.walletRepo.ListWithdrawalsByPhotographer(ctx, photographerID, 1, 1)
	for _, req := range existingReqs {
		if req.Status == domain.WithdrawalPending {
			return nil, errors.New("you already have a pending withdrawal request")
		}
	}

	// Use DB transaction to atomically deduct balance and create withdrawal request
	var withdrawalReq *domain.WithdrawalRequest
	err = u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Deduct from wallet
		if err := u.userRepo.UpdateWalletBalance(txCtx, photographerID, -input.Amount); err != nil {
			return err
		}

		newBalance := user.WalletBalance - input.Amount
		walletTx := &domain.WalletTransaction{
			ID:             uuid.New(),
			PhotographerID: photographerID,
			Type:           domain.TransactionDebit,
			Amount:         input.Amount,
			Description:    fmt.Sprintf("Penarikan dana ke %s - %s", input.BankName, input.AccountNumber),
			BalanceAfter:   newBalance,
		}
		if err := u.walletRepo.AddTransaction(txCtx, walletTx); err != nil {
			return err
		}

		withdrawalReq = &domain.WithdrawalRequest{
			ID:             uuid.New(),
			PhotographerID: photographerID,
			Amount:         input.Amount,
			BankName:       input.BankName,
			AccountNumber:  input.AccountNumber,
			AccountHolder:  input.AccountHolder,
			Status:         domain.WithdrawalPending,
		}

		return u.walletRepo.CreateWithdrawal(txCtx, withdrawalReq)
	})

	if err != nil {
		return nil, err
	}

	return withdrawalReq, nil
}

// ApproveWithdrawal — Super Admin only
func (u *WalletUsecase) ApproveWithdrawal(ctx context.Context, withdrawalID, adminID uuid.UUID, notes string) error {
	req, err := u.walletRepo.FindWithdrawalByID(ctx, withdrawalID)
	if err != nil {
		return errors.New("withdrawal request not found")
	}

	if req.Status != domain.WithdrawalPending {
		return errors.New("this withdrawal is not in pending status")
	}

	now := time.Now()
	req.ProcessedAt = &now
	return u.walletRepo.UpdateWithdrawalStatus(ctx, withdrawalID, domain.WithdrawalApproved, adminID, notes)
}

// RejectWithdrawal — Super Admin only; refunds balance
func (u *WalletUsecase) RejectWithdrawal(ctx context.Context, withdrawalID, adminID uuid.UUID, notes string) error {
	req, err := u.walletRepo.FindWithdrawalByID(ctx, withdrawalID)
	if err != nil {
		return errors.New("withdrawal request not found")
	}

	if req.Status != domain.WithdrawalPending {
		return errors.New("this withdrawal is not in pending status")
	}

	return u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Refund balance
		if err := u.userRepo.UpdateWalletBalance(txCtx, req.PhotographerID, req.Amount); err != nil {
			return err
		}

		balance, _ := u.walletRepo.GetBalance(txCtx, req.PhotographerID)
		walletTx := &domain.WalletTransaction{
			ID:             uuid.New(),
			PhotographerID: req.PhotographerID,
			Type:           domain.TransactionCredit,
			Amount:         req.Amount,
			Description:    fmt.Sprintf("Refund penarikan ditolak: %s", notes),
			BalanceAfter:   balance + req.Amount,
		}
		if err := u.walletRepo.AddTransaction(txCtx, walletTx); err != nil {
			return err
		}

		return u.walletRepo.UpdateWithdrawalStatus(txCtx, withdrawalID, domain.WithdrawalRejected, adminID, notes)
	})
}
