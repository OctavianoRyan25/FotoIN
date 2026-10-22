package infra

import (
	"context"
	"fmt"
	"log"

	"fotoin/backend/internal/config"
	"fotoin/backend/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(ctx context.Context, cfg config.DatabaseConfig) (*gorm.DB, error) {
	// Build DSN without password field when it is empty to prevent the
	// pg driver from falling back to the default 'postgres' database.
	var dsn string
	if cfg.Password != "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode,
		)
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			cfg.Host, cfg.User, cfg.Name, cfg.Port, cfg.SSLMode,
		)
	}

	log.Printf("Connecting to DB: host=%s user=%s dbname=%s port=%s sslmode=%s", cfg.Host, cfg.User, cfg.Name, cfg.Port, cfg.SSLMode)

	logLevel := logger.Silent
	if cfg.Environment == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("open database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database connection: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}

func AutoMigrate(ctx context.Context, db *gorm.DB) error {
	var dbName string
	if err := db.WithContext(ctx).Raw("SELECT current_database()").Scan(&dbName).Error; err != nil {
		return fmt.Errorf("get current database name: %w", err)
	}

	log.Println("Connected database:", dbName)

	err := db.WithContext(ctx).AutoMigrate(
		&model.UserModel{},
		&model.EventModel{},
		&model.GalleryModel{},
		&model.PhotoModel{},
		&model.OrderModel{},
		&model.OrderItemModel{},
		&model.WalletTransactionModel{},
		&model.WithdrawalRequestModel{},
		&model.SubscriptionModel{},
	)
	if err != nil {
		return fmt.Errorf("auto-migrate database: %w", err)
	}
	log.Println("Database migration completed")
	return nil
}
