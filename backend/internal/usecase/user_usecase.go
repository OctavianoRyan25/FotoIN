package usecase

import (
	"context"
	"errors"
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/util"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	userRepo domain.UserRepository
	mail     domain.MailService
}

// ── Input types ───────────────────────────────────────────────────────────────

type UpdateProfileInput struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
}

type UpdatePasswordInput struct {
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type CreatePasswordInput struct {
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ── Constructor ───────────────────────────────────────────────────────────────

func NewUserUsecase(userRepo domain.UserRepository, mail domain.MailService) *UserUsecase {
	return &UserUsecase{userRepo: userRepo, mail: mail}
}

// ── Methods ───────────────────────────────────────────────────────────────────

// EditProfile updates the authenticated user's display name.
func (u *UserUsecase) EditProfile(ctx context.Context, id uuid.UUID, input UpdateProfileInput) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	user.Name = input.Name
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, errors.New("gagal menyimpan profil")
	}
	return user, nil
}

// UpdatePassword verifies the old password, then stores the new bcrypt hash.
// Fix: was comparing plain-text old_password against the stored hash directly (always false).
func (u *UserUsecase) UpdatePassword(ctx context.Context, id uuid.UUID, input UpdatePasswordInput) error {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}
	if user.PasswordHash == "" {
		return errors.New("password belum diatur, silakan buat password terlebih dahulu")
	}

	// Verify old password against the stored bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.OldPassword)); err != nil {
		return errors.New("password lama tidak sesuai")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengenkripsi password")
	}
	return u.userRepo.UpdatePassword(ctx, id, string(newHash))
}

// CreatePassword sets a password for OAuth-only users who don't have one yet.
// Fix: was storing the plain-text password, not the bcrypt hash.
func (u *UserUsecase) CreatePassword(ctx context.Context, id uuid.UUID, input CreatePasswordInput) error {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}
	if user.PasswordHash != "" {
		return errors.New("password sudah diatur, silakan gunakan update password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengenkripsi password")
	}
	return u.userRepo.UpdatePassword(ctx, id, string(hash))
}

// DeleteAccount permanently removes the user's account.
func (u *UserUsecase) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	if _, err := u.userRepo.FindByID(ctx, id); err != nil {
		return errors.New("user tidak ditemukan")
	}
	return u.userRepo.Delete(ctx, id)
}

// ForgetPassword generates an OTP and sends it to the user's email.
func (u *UserUsecase) ForgetPassword(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("email tidak terdaftar")
	}

	otp := util.GenerateOTP()
	expiresAt := time.Now().Add(10 * time.Minute)

	if err := u.userRepo.AddOTP(ctx, otp, &expiresAt, user.ID); err != nil {
		return errors.New("gagal membuat OTP")
	}

	// Send email asynchronously so it doesn't block the API response
	go func() {
		_ = u.mail.SendOTP(user.Email, "Reset Password FotoIN", "Kode OTP Anda adalah: "+otp)
	}()

	return nil
}

type ResetPasswordInput struct {
	Email       string `json:"email" validate:"required,email"`
	OTP         string `json:"otp" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ResetPassword verifies the OTP and sets a new password.
func (u *UserUsecase) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	user, err := u.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return errors.New("email tidak terdaftar")
	}

	if user.OTPHash == "" || user.OTPHash != input.OTP {
		return errors.New("kode OTP salah")
	}

	if user.OTPExpiresAt == nil || user.OTPExpiresAt.Before(time.Now()) {
		return errors.New("kode OTP sudah kadaluarsa")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengenkripsi password baru")
	}

	// Update password and clear the OTP
	if err := u.userRepo.UpdatePassword(ctx, user.ID, string(newHash)); err != nil {
		return errors.New("gagal menyimpan password baru")
	}

	// Clear OTP so it can't be reused
	_ = u.userRepo.AddOTP(ctx, "", nil, user.ID)

	return nil
}
