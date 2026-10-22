package handler

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fotoin/backend/internal/usecase"
	"fotoin/backend/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderUsecase *usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUsecase: orderUsecase}
}

// CreateOrder — client-facing, no auth required
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var input usecase.CreateOrderInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	orderSummary, err := h.orderUsecase.CreateOrder(c.Request().Context(), input)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.Created(c, "Order created successfully", orderSummary)
}

// GetOrder — public, buyer must provide their email to verify ownership
func (h *OrderHandler) GetOrder(c echo.Context) error {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid order ID", nil)
	}

	email := c.QueryParam("email")
	if email == "" {
		return response.BadRequest(c, "Email is required", nil)
	}

	orderSummary, err := h.orderUsecase.GetOrder(c.Request().Context(), orderID, email)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.OK(c, "Order fetched", orderSummary)
}

// GetMyOrders — authenticated client; returns all orders by their JWT email
func (h *OrderHandler) GetMyOrders(c echo.Context) error {
	email, _ := c.Get("user_email").(string)
	if email == "" {
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

	orders, total, err := h.orderUsecase.GetOrdersByEmail(c.Request().Context(), email, page, limit)
	if err != nil {
		return response.InternalError(c, "Failed to fetch orders")
	}

	return response.Paginated(c, "Orders fetched", orders, page, limit, total)
}

// PaymentWebhook — receives notification from Midtrans or Duitku.
func (h *OrderHandler) PaymentWebhook(c echo.Context) error {
	var payload map[string]any
	if err := c.Bind(&payload); err != nil {
		return response.BadRequest(c, "Invalid webhook payload", err.Error())
	}

	if isDuidkuCallback(payload) {
		statusCode := fmt.Sprintf("%v", payload["statusCode"])
		if !isSuccessfulDuidkuStatus(statusCode) {
			return response.OK(c, "Webhook processed", nil)
		}

		orderID := extractOrderID(payload)
		gatewayID := extractString(payload, "reference")
		if gatewayID == "" {
			gatewayID = orderID
		}
		if orderID == "" {
			return response.BadRequest(c, "Invalid order ID", nil)
		}
		if err := h.orderUsecase.ProcessPaymentSuccess(c.Request().Context(), orderID, gatewayID, statusCode); err != nil {
			return response.InternalError(c, "Failed to process payment")
		}
		return response.OK(c, "Webhook processed", nil)
	}

	var midtransPayload struct {
		OrderID           string `json:"order_id"`
		TransactionID     string `json:"transaction_id"`
		TransactionStatus string `json:"transaction_status"`
		FraudStatus       string `json:"fraud_status"`
		GrossAmount       string `json:"gross_amount"`
	}
	if err := c.Bind(&midtransPayload); err != nil {
		return response.BadRequest(c, "Invalid webhook payload", err.Error())
	}

	if (midtransPayload.TransactionStatus == "capture" && midtransPayload.FraudStatus == "accept") ||
		midtransPayload.TransactionStatus == "settlement" {
		if err := h.orderUsecase.ProcessPaymentSuccess(
			c.Request().Context(),
			midtransPayload.OrderID,
			midtransPayload.TransactionID,
			midtransPayload.TransactionStatus,
		); err != nil {
			return response.InternalError(c, "Failed to process payment")
		}
	}

	return response.OK(c, "Webhook processed", nil)
}

func isDuidkuCallback(payload map[string]any) bool {
	_, hasMerchantOrderID := payload["merchantOrderId"]
	_, hasStatusCode := payload["statusCode"]
	_, hasMerchantOrderID2 := payload["merchant_order_id"]
	_, hasStatusCode2 := payload["status_code"]
	return hasMerchantOrderID || hasStatusCode || hasMerchantOrderID2 || hasStatusCode2
}

func isSuccessfulDuidkuStatus(statusCode string) bool {
	return statusCode == "00" || statusCode == "0000"
}

func extractOrderID(payload map[string]any) string {
	for _, key := range []string{"merchantOrderId", "merchant_order_id", "order_id", "orderId"} {
		if value, ok := payload[key]; ok {
			if str, ok := value.(string); ok && str != "" {
				return str
			}
		}
	}
	return ""
}

func extractString(payload map[string]any, key string) string {
	if value, ok := payload[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
		if b, err := json.Marshal(value); err == nil {
			return string(b)
		}
	}
	return ""
}


