package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TransactionType string
type WithdrawalStatus string

const (
	TransactionCredit TransactionType = "credit"
	TransactionDebit  TransactionType = "debit"

	WithdrawalPending  WithdrawalStatus = "pending"
	WithdrawalApproved WithdrawalStatus = "approved"
	WithdrawalRejected WithdrawalStatus = "rejected"
)

type WalletTransaction struct {
	ID             uuid.UUID       `json:"id"`
	PhotographerID uuid.UUID       `json:"photographer_id"`
	OrderID        *uuid.UUID      `json:"order_id,omitempty"`
	Type           TransactionType `json:"type"`
	Amount         float64         `json:"amount"`
	Description    string          `json:"description"`
	BalanceAfter   float64         `json:"balance_after"`
	CreatedAt      time.Time       `json:"created_at"`
}

type WithdrawalRequest struct {
	ID             uuid.UUID        `json:"id"`
	PhotographerID uuid.UUID        `json:"photographer_id"`
	Amount         float64          `json:"amount"`
	BankName       string           `json:"bank_name"`
	AccountNumber  string           `json:"account_number"`
	AccountHolder  string           `json:"account_holder"`
	Status         WithdrawalStatus `json:"status"`
	ProcessedBy    *uuid.UUID       `json:"processed_by,omitempty"`
	ProcessedAt    *time.Time       `json:"processed_at,omitempty"`
	Notes          *string          `json:"notes,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`

	Photographer *User `json:"photographer,omitempty"`
	Processor    *User `json:"processor,omitempty"`
}

type WalletRepository interface {
	GetBalance(ctx context.Context, photographerID uuid.UUID) (float64, error)
	AddTransaction(ctx context.Context, tx *WalletTransaction) error
	GetTransactions(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]WalletTransaction, int64, error)
	CreateWithdrawal(ctx context.Context, req *WithdrawalRequest) error
	FindWithdrawalByID(ctx context.Context, id uuid.UUID) (*WithdrawalRequest, error)
	ListWithdrawals(ctx context.Context, status WithdrawalStatus, page, limit int) ([]WithdrawalRequest, int64, error)
	ListWithdrawalsByPhotographer(ctx context.Context, photographerID uuid.UUID, page, limit int) ([]WithdrawalRequest, int64, error)
	UpdateWithdrawalStatus(ctx context.Context, id uuid.UUID, status WithdrawalStatus, processedBy uuid.UUID, notes string) error
}
