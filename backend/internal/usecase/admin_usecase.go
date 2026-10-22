package usecase

import (
	"context"
	"fotoin/backend/internal/domain"

	"github.com/google/uuid"
)

type AdminUsecase struct {
	adminRepo domain.AdminRepository
	userRepo  domain.UserRepository
}

func NewAdminUsecase(adminRepo domain.AdminRepository, userRepo domain.UserRepository) *AdminUsecase {
	return &AdminUsecase{
		adminRepo: adminRepo,
		userRepo:  userRepo,
	}
}

func (u *AdminUsecase) GetRevenueStats(ctx context.Context) (*domain.AdminRevenueStats, error) {
	return u.adminRepo.GetRevenueStats(ctx)
}

func (u *AdminUsecase) ListUsers(ctx context.Context, page, limit int, search, role string) ([]domain.User, int64, error) {
	return u.userRepo.List(ctx, page, limit, search, role)
}

func (u *AdminUsecase) SuspendUser(ctx context.Context, userID uuid.UUID, suspended bool) error {
	return u.userRepo.SuspendUser(ctx, userID, suspended)
}

func (u *AdminUsecase) GetUserProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.HasPassword = user.PasswordHash != ""
	return user, nil
}
