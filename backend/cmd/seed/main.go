// cmd/seed/main.go — Seeds the database with:
//   1. A Super Admin account
//   2. A Photographer account (for testing)
//
// Usage (from backend/ directory):
//   go run ./cmd/seed/main.go
//
// Environment: reads the same .env as the API server.

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"fotoin/backend/internal/config"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"
	"fotoin/backend/internal/repository"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type seedUser struct {
	Name     string
	Email    string
	Password string
	Role     domain.UserRole
	Tier     domain.SubscriptionTier
}

var seeds = []seedUser{
	{
		Name:     "Super Admin FotoIN",
		Email:    "admin@fotoin.id",
		Password: "Admin@123!",
		Role:     domain.RoleSuperAdmin,
		Tier:     domain.TierPremium,
	},
	{
		Name:     "Rangga Foto",
		Email:    "fotografer@fotoin.id",
		Password: "Foto@123!",
		Role:     domain.RolePhotographer,
		Tier:     domain.TierPremium,
	},
}

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	if err := run(); err != nil {
		log.Printf("Seeder failed: %v", err)
		os.Exit(1)
	}
}

func run() (err error) {
	ctx := context.Background()
	cfg := config.LoadConfig()
	db, err := infra.NewDatabase(ctx, cfg.Database)
	if err != nil {
		return err
	}
	defer func() {
		sqlDB, closeErr := db.DB()
		if closeErr == nil {
			closeErr = sqlDB.Close()
		}
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close database: %w", closeErr))
		}
	}()

	// Auto-migrate so tables exist
	if err := infra.AutoMigrate(ctx, db); err != nil {
		return err
	}

	log.Println("─── FotoIN Database Seeder ───────────────────────────────")

	userRepo := repository.NewUserRepository(db)
	for _, s := range seeds {
		upsertUser(userRepo, s)
	}

	log.Println("─── Seeding complete ─────────────────────────────────────")
	printSummary()
	return nil
}

func upsertUser(userRepo domain.UserRepository, s seedUser) {
	ctx := context.Background()
	existing, err := userRepo.FindByEmail(ctx, s.Email)

	if err == nil && existing != nil {
		log.Printf("  [SKIP] %s (%s) already exists\n", s.Name, s.Email)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("  [ERROR] bcrypt for %s: %v", s.Email, err)
	}

	now := time.Now().Add(365 * 24 * time.Hour) // 1-year premium for test accounts
	user := &domain.User{
		ID:                    uuid.New(),
		Name:                  s.Name,
		Email:                 s.Email,
		PasswordHash:          string(hash),
		Role:                  s.Role,
		SubscriptionTier:      s.Tier,
		SubscriptionExpiresAt: &now,
		IsVerified:            true,
		IsSuspended:           false,
	}

	if err := userRepo.Create(ctx, user); err != nil {
		log.Printf("  [ERROR] Creating %s: %v\n", s.Email, err)
		return
	}

	roleLabel := string(s.Role)
	log.Printf("  [OK]   %-20s %-35s role=%-15s tier=%s\n", s.Name, s.Email, roleLabel, s.Tier)
}

func printSummary() {
	separator := "─────────────────────────────────────────────────────────"
	fmt.Println(separator)
	fmt.Println("  Credential Summary (save these!)")
	fmt.Println(separator)
	for _, s := range seeds {
		fmt.Printf("  %-20s  email: %-30s  pass: %s\n", "["+string(s.Role)+"]", s.Email, s.Password)
	}
	fmt.Println(separator)
	fmt.Println()
	fmt.Printf("  Dashboard: http://localhost:5173  (login with %s)\n", os.Getenv("APP_PORT"))
	fmt.Println("  Client:    http://localhost:3000")
	fmt.Println(separator)
}
