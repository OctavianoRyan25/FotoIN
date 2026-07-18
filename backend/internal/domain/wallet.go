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
	ID             uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PhotographerID uuid.UUID       `gorm:"type:uuid;not null;index" json:"photographer_id"`
	OrderID        *uuid.UUID      `gorm:"type:uuid;null" json:"order_id,omitempty"`
	Type           TransactionType `gorm:"type:varchar(10);not null" json:"type"`
	Amount         float64         `gorm:"type:decimal(14,2);not null" json:"amount"`
	Description    string          `gorm:"size:500" json:"description"`
	BalanceAfter   float64         `gorm:"type:decimal(14,2);not null" json:"balance_after"`
	CreatedAt      time.Time       `json:"created_at"`
}

type WithdrawalRequest struct {
	ID              uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PhotographerID  uuid.UUID        `gorm:"type:uuid;not null;index" json:"photographer_id"`
	Amount          float64          `gorm:"type:decimal(14,2);not null" json:"amount"`
	BankName        string           `gorm:"size:100;not null" json:"bank_name"`
	AccountNumber   string           `gorm:"size:50;not null" json:"account_number"`
	AccountHolder   string           `gorm:"size:255;not null" json:"account_holder"`
	Status          WithdrawalStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	ProcessedBy     *uuid.UUID       `gorm:"type:uuid;null" json:"processed_by,omitempty"`
	ProcessedAt     *time.Time       `gorm:"null" json:"processed_at,omitempty"`
	Notes           *string          `gorm:"type:text;null" json:"notes,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`

	// Relations
	Photographer User  `gorm:"foreignKey:PhotographerID" json:"photographer,omitempty"`
	Processor    *User `gorm:"foreignKey:ProcessedBy" json:"processor,omitempty"`
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
