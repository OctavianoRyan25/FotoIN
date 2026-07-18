package infra

import (
	"fmt"
	"log"
	"os"

	"fotoin/backend/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	// Build DSN without password field when it is empty to prevent the
	// pg driver from falling back to the default 'postgres' database.
	var dsn string
	if password != "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			host, user, password, dbname, port, sslmode,
		)
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			host, user, dbname, port, sslmode,
		)
	}

	log.Printf("Connecting to DB: host=%s user=%s dbname=%s port=%s sslmode=%s", host, user, dbname, port, sslmode)

	logLevel := logger.Silent
	if os.Getenv("APP_ENV") == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
	return db
}

func AutoMigrate(db *gorm.DB) {
	var dbName string
	db.Raw("SELECT current_database()").Scan(&dbName)

	log.Println("Connected database:", dbName)

	err := db.AutoMigrate(
		&domain.User{},
		&domain.Event{},
		&domain.Gallery{},
		&domain.Photo{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.WalletTransaction{},
		&domain.WithdrawalRequest{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migration completed")
}
