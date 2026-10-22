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

type WalletHandler struct {
	walletUsecase *usecase.WalletUsecase
}

func NewWalletHandler(walletUsecase *usecase.WalletUsecase) *WalletHandler {
	return &WalletHandler{walletUsecase: walletUsecase}
}

func (h *WalletHandler) GetWallet(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	summary, err := h.walletUsecase.GetSummary(c.Request().Context(), photographerID, page, limit)
	if err != nil {
		return response.InternalError(c, "Failed to fetch wallet data")
	}

	return response.OK(c, "Wallet fetched", summary)
}

func (h *WalletHandler) RequestWithdrawal(c echo.Context) error {
	photographerID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	var input usecase.WithdrawalInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	req, err := h.walletUsecase.RequestWithdrawal(c.Request().Context(), photographerID, input)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "Withdrawal request submitted", req)
}

// Admin handlers for withdrawals
func (h *WalletHandler) ListWithdrawals(c echo.Context) error {
	status := domain.WithdrawalStatus(c.QueryParam("status"))
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	reqs, total, err := h.walletUsecase.ListWithdrawals(c.Request().Context(), status, page, limit)
	if err != nil {
		return response.InternalError(c, "Failed to fetch withdrawals")
	}

	return response.Paginated(c, "Withdrawals fetched", reqs, page, limit, total)
}

func (h *WalletHandler) ApproveWithdrawal(c echo.Context) error {
	adminID, err := middleware.GetUserID(c)
	if err != nil {
		return response.Unauthorized(c, "Invalid token")
	}

	withdrawalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid withdrawal ID", nil)
	}

	var body struct {
		Action string `json:"action"` // "approve" or "reject"
		Notes  string `json:"notes"`
	}
	if err := c.Bind(&body); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	switch body.Action {
	case "approve":
		err = h.walletUsecase.ApproveWithdrawal(c.Request().Context(), withdrawalID, adminID, body.Notes)
	case "reject":
		err = h.walletUsecase.RejectWithdrawal(c.Request().Context(), withdrawalID, adminID, body.Notes)
	default:
		return response.BadRequest(c, "Action must be 'approve' or 'reject'", nil)
	}

	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.OK(c, "Withdrawal "+body.Action+"d successfully", nil)
}
