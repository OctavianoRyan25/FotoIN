package handler

import (
	"strconv"

	"fotoin/backend/internal/domain"
	"fotoin/backend/internal/middleware"
	"fotoin/backend/internal/usecase"
	"fotoin/backend/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SubscriptionHandler struct {
	subUsecase *usecase.SubscriptionUsecase
}

func NewSubscriptionHandler(subUsecase *usecase.SubscriptionUsecase) *SubscriptionHandler {
	return &SubscriptionHandler{subUsecase: subUsecase}
}

func (h *SubscriptionHandler) GetMySubscription(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	res, err := h.subUsecase.GetMySubscription(c.Request().Context(), userID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, "Subscription status fetched", res)
}

func (h *SubscriptionHandler) Upgrade(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	res, err := h.subUsecase.UpgradeToPremium(c.Request().Context(), userID)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "Upgrade subscription created", res)
}

func (h *SubscriptionHandler) SimulatePayment(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	subID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid subscription ID", nil)
	}

	role, _ := middleware.GetUserRole(c)
	isAdmin := role == domain.RoleSuperAdmin

	sub, err := h.subUsecase.SimulatePayment(c.Request().Context(), subID, userID, isAdmin)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.OK(c, "Subscription payment simulated successfully! Your account is now Premium.", sub)
}

func (h *SubscriptionHandler) ListAdminSubscriptions(c echo.Context) error {
	status := domain.SubscriptionStatus(c.QueryParam("status"))
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	subs, total, err := h.subUsecase.ListAllSubscriptions(c.Request().Context(), status, page, limit)
	if err != nil {
		return response.InternalError(c, "Failed to fetch subscriptions")
	}

	return response.Paginated(c, "Subscriptions fetched", subs, page, limit, total)
}

func (h *SubscriptionHandler) AdminGrantPremium(c echo.Context) error {
	var body struct {
		UserID string `json:"user_id" validate:"required"`
		Months int    `json:"months"`
	}
	if err := c.Bind(&body); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	targetUserID, err := uuid.Parse(body.UserID)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", nil)
	}

	user, err := h.subUsecase.AdminGrantPremium(c.Request().Context(), targetUserID, body.Months)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.OK(c, "Premium tier successfully granted to user", user)
}
