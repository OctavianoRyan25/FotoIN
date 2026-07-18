package handler

import (
	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/middleware"
	"fotoin/backend/pkg/response"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AdminHandler struct {
	userRepo domain.UserRepository
	db       *gorm.DB
	redis    *redis.Client
}

func NewAdminHandler(userRepo domain.UserRepository, db *gorm.DB, redis *redis.Client) *AdminHandler {
	return &AdminHandler{userRepo: userRepo, db: db, redis: redis}
}

func (h *AdminHandler) ListUsers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	ctx := c.Request().Context()
	users, total, err := h.userRepo.List(ctx, page, limit, search)
	if err != nil {
		return response.InternalError(c, "Failed to fetch users")
	}

	return response.Paginated(c, "Users fetched", users, page, limit, total)
}

func (h *AdminHandler) SuspendUser(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	var body struct {
		Suspended bool `json:"suspended"`
	}
	if err := c.Bind(&body); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	ctx := c.Request().Context()
	if err := h.userRepo.SuspendUser(ctx, userID, body.Suspended); err != nil {
		return response.InternalError(c, "Failed to update user status")
	}

	msg := "User suspended"
	if !body.Suspended {
		msg = "User unsuspended"
	}
	return response.OK(c, msg, nil)
}

func (h *AdminHandler) GetRevenue(c echo.Context) error {
	var result struct {
		TotalRevenue       float64 `json:"total_revenue"`
		BasicTierRevenue   float64 `json:"basic_tier_revenue"`
		PremiumTierRevenue float64 `json:"premium_tier_revenue"`
		TotalOrders        int64   `json:"total_orders"`
		TotalPhotographers int64   `json:"total_photographers"`
	}

	h.db.Model(&domain.Order{}).
		Select("COALESCE(SUM(platform_fee), 0) as total_revenue, COUNT(*) as total_orders").
		Where("payment_status = ?", domain.PaymentPaid).
		Scan(&result)

	h.db.Model(&domain.User{}).
		Where("role = ?", domain.RolePhotographer).
		Count(&result.TotalPhotographers)

	return response.OK(c, "Revenue overview fetched", result)
}

func (h *AdminHandler) GetMe(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	ctx := c.Request().Context()
	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	user.HasPassword = user.PasswordHash != ""

	return response.OK(c, "Profile fetched", user)
}
