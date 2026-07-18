package infra

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"fotoin/backend/internal/domain"
)

// MidtransQRISService calls the Midtrans Core API (NOT Snap) to generate QRIS.
// Supported payment_type values: "qris" (universal QRIS), "gopay" (GoPay QR).
type MidtransQRISService struct {
	serverKey    string
	isProduction bool
}

func NewMidtransQRISService() *MidtransQRISService {
	env := os.Getenv("MIDTRANS_ENV")
	return &MidtransQRISService{
		serverKey:    os.Getenv("MIDTRANS_SERVER_KEY"),
		isProduction: env == "production",
	}
}

func (s *MidtransQRISService) baseURL() string {
	if s.isProduction {
		return "https://api.midtrans.com"
	}
	return "https://api.sandbox.midtrans.com"
}

// basicAuth encodes "SERVER_KEY:" as Base64 per Midtrans docs.
func (s *MidtransQRISService) basicAuth() string {
	return base64.StdEncoding.EncodeToString([]byte(s.serverKey + ":"))
}

// GenerateQRCode creates a QRIS charge via Midtrans Core API /v2/charge.
// The response contains:
//   - qr_string → raw QRIS string (scan with banking/e-wallet app)
//   - actions[].url where name=="generate-qr-code" → hosted QR image URL
func (s *MidtransQRISService) GenerateQRCode(req domain.QRCodePaymentRequest) (*domain.QRCodePaymentResponse, error) {
	// Build request payload — always use payment_type: "qris"
	payload := map[string]interface{}{
		"payment_type": "qris",
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": int64(req.Amount),
		},
		"qris": map[string]interface{}{
			"acquirer": "gopay", // gopay acquirer supports most banks & e-wallets
		},
		"customer_details": map[string]interface{}{
			"email": req.CustomerEmail,
		},
		"item_details": []map[string]interface{}{
			{
				"id":       req.OrderID,
				"price":    int64(req.Amount),
				"quantity": 1,
				"name":     req.Description,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/v2/charge", s.baseURL())
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Basic "+s.basicAuth())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("midtrans request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("midtrans error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse Midtrans response
	var result struct {
		StatusCode        string `json:"status_code"`
		StatusMessage     string `json:"status_message"`
		TransactionID     string `json:"transaction_id"`
		OrderID           string `json:"order_id"`
		QRString          string `json:"qr_string"`   // raw QRIS string (scan with any app)
		ExpiryTime        string `json:"expiry_time"` // "2024-07-09 15:30:00"
		Actions           []struct {
			Name   string `json:"name"`
			Method string `json:"method"`
			URL    string `json:"url"`
		} `json:"actions"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	// qr_string is the QRIS EMV string — use it with a QR rendering lib on the frontend
	qrString := result.QRString

	// Find hosted QR image URL from actions
	qrImageURL := ""
	for _, action := range result.Actions {
		if action.Name == "generate-qr-code" {
			qrImageURL = action.URL
			break
		}
	}

	// For QRIS, at minimum qr_string should be present
	if qrString == "" && qrImageURL == "" {
		return nil, errors.New("midtrans did not return a QR code in response")
	}

	return &domain.QRCodePaymentResponse{
		OrderID:      req.OrderID,
		Gateway:      domain.PaymentGatewayMidtrans,
		QRCodeURL:    qrImageURL,   // hosted image (may be empty in sandbox)
		QRCodeString: qrString,     // EMV string → render on frontend
		ExpiresAt:    result.ExpiryTime,
		Amount:       req.Amount,
		PaymentMethod: domain.PaymentMethodQRIS,
	}, nil
}

// CheckPaymentStatus fetches the current transaction status from Midtrans.
// Returns Midtrans transaction_status: "pending", "settlement", "deny", "expire", etc.
func (s *MidtransQRISService) CheckPaymentStatus(orderID string) (string, error) {
	url := fmt.Sprintf("%s/v2/%s/status", s.baseURL(), orderID)

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Basic "+s.basicAuth())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("midtrans status check: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		TransactionStatus string `json:"transaction_status"`
		FraudStatus       string `json:"fraud_status"`
		StatusMessage     string `json:"status_message"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parse status: %w", err)
	}
	return result.TransactionStatus, nil
}
