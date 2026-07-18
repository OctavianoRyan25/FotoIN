package infra

import (
	"context"
	"fotoin/backend/internal/domain"

	"gorm.io/gorm"
)

type txKey struct{}

type gormTxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) domain.TxManager {
	return &gormTxManager{db: db}
}

func (m *gormTxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxWithTx := context.WithValue(ctx, txKey{}, tx)
		return fn(ctxWithTx)
	})
}

// Dipakai semua repo untuk ambil *gorm.DB yang tepat (tx kalau ada, db biasa kalau tidak)
func GetExecutor(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return db.WithContext(ctx)
}
