package domain

// PaymentMethod represents the QRIS payment method.
// FotoIN only supports QRIS via Midtrans Core API (no Snap).
type PaymentMethod string

const (
	PaymentMethodQRIS PaymentMethod = "qris"
)

// PaymentGateway identifies the gateway provider.
type PaymentGateway string

const (
	PaymentGatewayMidtrans PaymentGateway = "midtrans"
	PaymentGatewayDuidku   PaymentGateway = "duitku"
)

// QRCodePaymentRequest is sent to the payment service to generate a QR code.
type QRCodePaymentRequest struct {
	OrderID       string        `json:"order_id"`
	Amount        float64       `json:"amount"`
	CustomerEmail string        `json:"customer_email"`
	Description   string        `json:"description"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

// QRCodePaymentResponse contains QR data returned to the client.
// QRCodeString is the raw QRIS EMV string — use a QR renderer on the frontend.
// QRCodeURL is the hosted QR image (may be empty in sandbox).
type QRCodePaymentResponse struct {
	OrderID       string         `json:"order_id"`
	Gateway       PaymentGateway `json:"gateway"`
	QRCodeURL     string         `json:"qr_code_url"`    // hosted image URL (optional)
	QRCodeString  string         `json:"qr_code_string"` // raw EMV string → render as QR
	PaymentURL    string         `json:"payment_url,omitempty"`
	ExpiresAt     string         `json:"expires_at"`
	Amount        float64        `json:"amount"`
	PaymentMethod PaymentMethod  `json:"payment_method"`
}

// PaymentService is the port for payment gateway operations.
type PaymentService interface {
	// GenerateQRCode creates a QRIS charge and returns QR data.
	GenerateQRCode(req QRCodePaymentRequest) (*QRCodePaymentResponse, error)
	// CheckPaymentStatus queries the current transaction status by order ID.
	CheckPaymentStatus(orderID string) (string, error)
}
