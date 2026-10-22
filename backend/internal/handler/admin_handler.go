package handler

import (
	"strconv"

	"fotoin/backend/internal/middleware"
	"fotoin/backend/internal/usecase"
	"fotoin/backend/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	adminUsecase *usecase.AdminUsecase
}

func NewAdminHandler(adminUsecase *usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{adminUsecase: adminUsecase}
}

func (h *AdminHandler) ListUsers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")
	role := c.QueryParam("role")
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	ctx := c.Request().Context()
	users, total, err := h.adminUsecase.ListUsers(ctx, page, limit, search, role)
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
	if err := h.adminUsecase.SuspendUser(ctx, userID, body.Suspended); err != nil {
		return response.InternalError(c, "Failed to update user status")
	}

	msg := "User suspended"
	if !body.Suspended {
		msg = "User unsuspended"
	}
	return response.OK(c, msg, nil)
}

func (h *AdminHandler) GetRevenue(c echo.Context) error {
	ctx := c.Request().Context()
	stats, err := h.adminUsecase.GetRevenueStats(ctx)
	if err != nil {
		return response.InternalError(c, "Failed to fetch revenue overview")
	}

	return response.OK(c, "Revenue overview fetched", stats)
}

func (h *AdminHandler) GetMe(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	ctx := c.Request().Context()
	user, err := h.adminUsecase.GetUserProfile(ctx, userID)
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	return response.OK(c, "Profile fetched", user)
}
