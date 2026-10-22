package model

import (
	"time"

	"fotoin/backend/internal/domain"

	"github.com/google/uuid"
)

type WalletTransactionModel struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PhotographerID uuid.UUID  `gorm:"type:uuid;not null;index"`
	OrderID        *uuid.UUID `gorm:"type:uuid;null"`
	Type           string     `gorm:"type:varchar(10);not null"`
	Amount         float64    `gorm:"type:decimal(14,2);not null"`
	Description    string     `gorm:"size:500"`
	BalanceAfter   float64    `gorm:"type:decimal(14,2);not null"`
	CreatedAt      time.Time
}

func (WalletTransactionModel) TableName() string {
	return "wallet_transactions"
}

func (m *WalletTransactionModel) ToDomain() *domain.WalletTransaction {
	if m == nil {
		return nil
	}
	return &domain.WalletTransaction{
		ID:             m.ID,
		PhotographerID: m.PhotographerID,
		OrderID:        m.OrderID,
		Type:           domain.TransactionType(m.Type),
		Amount:         m.Amount,
		Description:    m.Description,
		BalanceAfter:   m.BalanceAfter,
		CreatedAt:      m.CreatedAt,
	}
}

func FromDomainWalletTransaction(tx *domain.WalletTransaction) *WalletTransactionModel {
	if tx == nil {
		return nil
	}
	return &WalletTransactionModel{
		ID:             tx.ID,
		PhotographerID: tx.PhotographerID,
		OrderID:        tx.OrderID,
		Type:           string(tx.Type),
		Amount:         tx.Amount,
		Description:    tx.Description,
		BalanceAfter:   tx.BalanceAfter,
		CreatedAt:      tx.CreatedAt,
	}
}

func ToDomainWalletTransactionList(models []WalletTransactionModel) []domain.WalletTransaction {
	res := make([]domain.WalletTransaction, len(models))
	for i, m := range models {
		res[i] = *m.ToDomain()
	}
	return res
}

type WithdrawalRequestModel struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PhotographerID uuid.UUID  `gorm:"type:uuid;not null;index"`
	Amount         float64    `gorm:"type:decimal(14,2);not null"`
	BankName       string     `gorm:"size:100;not null"`
	AccountNumber  string     `gorm:"size:50;not null"`
	AccountHolder  string     `gorm:"size:255;not null"`
	Status         string     `gorm:"type:varchar(20);not null;default:'pending'"`
	ProcessedBy    *uuid.UUID `gorm:"type:uuid;null"`
	ProcessedAt    *time.Time `gorm:"null"`
	Notes          *string    `gorm:"type:text;null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// Relations
	Photographer UserModel   `gorm:"foreignKey:PhotographerID"`
	Processor    *UserModel  `gorm:"foreignKey:ProcessedBy"`
}

func (WithdrawalRequestModel) TableName() string {
	return "withdrawal_requests"
}

func (m *WithdrawalRequestModel) ToDomain() *domain.WithdrawalRequest {
	if m == nil {
		return nil
	}
	req := &domain.WithdrawalRequest{
		ID:             m.ID,
		PhotographerID: m.PhotographerID,
		Amount:         m.Amount,
		BankName:       m.BankName,
		AccountNumber:  m.AccountNumber,
		AccountHolder:  m.AccountHolder,
		Status:         domain.WithdrawalStatus(m.Status),
		ProcessedBy:    m.ProcessedBy,
		ProcessedAt:    m.ProcessedAt,
		Notes:          m.Notes,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}

	if m.Photographer.ID != uuid.Nil {
		req.Photographer = m.Photographer.ToDomain()
	}
	if m.Processor != nil && m.Processor.ID != uuid.Nil {
		req.Processor = m.Processor.ToDomain()
	}

	return req
}

func FromDomainWithdrawalRequest(req *domain.WithdrawalRequest) *WithdrawalRequestModel {
	if req == nil {
		return nil
	}
	return &WithdrawalRequestModel{
		ID:             req.ID,
		PhotographerID: req.PhotographerID,
		Amount:         req.Amount,
		BankName:       req.BankName,
		AccountNumber:  req.AccountNumber,
		AccountHolder:  req.AccountHolder,
		Status:         string(req.Status),
		ProcessedBy:    req.ProcessedBy,
		ProcessedAt:    req.ProcessedAt,
		Notes:          req.Notes,
		CreatedAt:      req.CreatedAt,
		UpdatedAt:      req.UpdatedAt,
	}
}

func ToDomainWithdrawalRequestList(models []WithdrawalRequestModel) []domain.WithdrawalRequest {
	res := make([]domain.WithdrawalRequest, len(models))
	for i, m := range models {
		res[i] = *m.ToDomain()
	}
	return res
}
