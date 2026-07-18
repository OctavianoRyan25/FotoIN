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
	"fmt"
	"log"
	"os"
	"time"

	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/infra"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

	db := infra.NewDatabase()

	// Auto-migrate so tables exist
	infra.AutoMigrate(db)

	log.Println("─── FotoIN Database Seeder ───────────────────────────────")

	for _, s := range seeds {
		upsertUser(db, s)
	}

	log.Println("─── Seeding complete ─────────────────────────────────────")
	printSummary()
}

func upsertUser(db *gorm.DB, s seedUser) {
	var existing domain.User
	err := db.Where("email = ?", s.Email).First(&existing).Error

	if err == nil {
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

	if err := db.Create(user).Error; err != nil {
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
