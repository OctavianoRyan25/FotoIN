package infra

import (
	"bytes"
	"context"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"fotoin/backend/internal/config"
	"fotoin/backend/internal/domain"
)

// ─── SNAP QRIS MPM DTOs ───────────────────────────────────────────────────────

type TokenResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	AccessToken     string `json:"accessToken"`
	TokenType       string `json:"tokenType"`
	ExpiresIn       string `json:"expiresIn"`
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type AdditionalInfo struct {
	Acquirer   string `json:"acquirer,omitempty"`
	MerchantID string `json:"merchantId,omitempty"`
}

type QRISRequest struct {
	PartnerReferenceNo string         `json:"partnerReferenceNo"`
	Amount             Amount         `json:"amount"`
	MerchantID         string         `json:"merchantId,omitempty"`
	ValidityPeriod     string         `json:"validityPeriod,omitempty"`
	AdditionalInfo     AdditionalInfo `json:"additionalInfo,omitempty"`
}

type QRISResponse struct {
	ResponseCode       string `json:"responseCode"`
	ResponseMessage    string `json:"responseMessage"`
	ReferenceNo        string `json:"referenceNo"`
	PartnerReferenceNo string `json:"partnerReferenceNo"`
	QRContent          string `json:"qrContent"`
	QRImage            string `json:"qrImage,omitempty"`
	QRURL              string `json:"qrUrl,omitempty"`
}

type QRISStatusQueryRequest struct {
	OriginalPartnerReferenceNo string         `json:"originalPartnerReferenceNo"`
	OriginalReferenceNo        string         `json:"originalReferenceNo,omitempty"`
	ServiceCode                string         `json:"serviceCode,omitempty"`
	Amount                     Amount         `json:"amount"`
	AdditionalInfo             AdditionalInfo `json:"additionalInfo,omitempty"`
}

type QRISStatusQueryResponse struct {
	ResponseCode               string `json:"responseCode"`
	ResponseMessage            string `json:"responseMessage"`
	OriginalPartnerReferenceNo string `json:"originalPartnerReferenceNo"`
	LatestTransactionStatus    string `json:"latestTransactionStatus"`
	TransactionStatusDesc      string `json:"transactionStatusDesc,omitempty"`
}

// ─── Service Definition ───────────────────────────────────────────────────────

// MidtransQRISService handles QRIS generation via Midtrans SNAP MPM QRIS (ASPI/BI standard)
// with backward-compatible fallback to Midtrans Core API /v2/charge.
type MidtransQRISService struct {
	serverKey    string
	isProduction bool
	clientId     string
	clientSecret string
	partnerId    string
	channelId    string
	merchantId   string
	privateKey   string
	client       *http.Client

	// In-memory token cache to avoid requesting B2B tokens on every transaction
	tokenMu     sync.RWMutex
	cachedToken string
	tokenExpiry time.Time
}

func NewMidtransQRISService(cfg config.MidtransConfig) *MidtransQRISService {
	channelId := cfg.ChannelId
	if channelId == "" {
		channelId = "95221" // Default standard ASPI QRIS MPM channel ID
	}

	partnerId := cfg.PartnerId
	if partnerId == "" {
		partnerId = cfg.MerchantId
	}

	return &MidtransQRISService{
		serverKey:    cfg.ServerKey,
		isProduction: cfg.Environment == "production",
		clientId:     cfg.ClientId,
		clientSecret: cfg.ClientSecret,
		partnerId:    partnerId,
		channelId:    channelId,
		merchantId:   cfg.MerchantId,
		privateKey:   cfg.PrivateKey,
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

// ─── Helper Methods ───────────────────────────────────────────────────────────

func (s *MidtransQRISService) baseURL() string {
	if s.isProduction {
		return "https://api.midtrans.com"
	}
	return "https://api.sandbox.midtrans.com"
}

func (s *MidtransQRISService) basicAuth() string {
	return base64.StdEncoding.EncodeToString([]byte(s.serverKey + ":"))
}

func (s *MidtransQRISService) timestamp() string {
	// ISO-8601 with WIB (+07:00) timezone offset
	return time.Now().In(time.FixedZone("WIB", 7*60*60)).Format("2006-01-02T15:04:05-07:00")
}

// loadRSAPrivateKey parses an RSA private key from either a file path or an inline PEM string.
// Supports both PKCS#8 (required by SNAP) and PKCS#1.
func (s *MidtransQRISService) loadRSAPrivateKey() (*rsa.PrivateKey, error) {
	if s.privateKey == "" {
		return nil, errors.New("midtrans private key tidak dikonfigurasi")
	}

	var data []byte
	var err error

	// If the privateKey string is a valid file path, read from file; otherwise treat as PEM string.
	if _, statErr := os.Stat(s.privateKey); statErr == nil {
		data, err = os.ReadFile(s.privateKey)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca private key file: %w", err)
		}
	} else {
		data = []byte(s.privateKey)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("format PEM private key tidak valid")
	}

	// Try PKCS#8 parsing first (SNAP standard)
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
	}

	// Fallback to PKCS#1
	if rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return rsaKey, nil
	}

	return nil, errors.New("private key harus berupa format RSA PKCS#8 atau PKCS#1")
}

// getAccessToken retrieves a valid B2B OAuth token, caching it until expiration.
func (s *MidtransQRISService) getAccessToken(ctx context.Context) (string, error) {
	s.tokenMu.RLock()
	if s.cachedToken != "" && time.Now().Before(s.tokenExpiry) {
		token := s.cachedToken
		s.tokenMu.RUnlock()
		return token, nil
	}
	s.tokenMu.RUnlock()

	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()

	// Double check inside write lock
	if s.cachedToken != "" && time.Now().Before(s.tokenExpiry) {
		return s.cachedToken, nil
	}

	ts := s.timestamp()
	privateKey, err := s.loadRSAPrivateKey()
	if err != nil {
		return "", err
	}

	// SNAP StringToSign for access token: ClientID + "|" + X-TIMESTAMP
	stringToSign := s.clientId + "|" + ts
	digest := sha256.Sum256([]byte(stringToSign))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest[:])
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	body, err := json.Marshal(map[string]string{
		"grantType": "client_credentials",
	})
	if err != nil {
		return "", err
	}

	endpoint := s.baseURL() + "/v1.0/access-token/b2b"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-TIMESTAMP", ts)
	req.Header.Set("X-CLIENT-KEY", s.clientId)
	req.Header.Set("X-SIGNATURE", base64.StdEncoding.EncodeToString(signature))

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request access token: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Midtrans access token HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var result TokenResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}
	if result.AccessToken == "" {
		return "", fmt.Errorf("access token kosong, response: %s", string(respBody))
	}

	// Cache token (default 900s, subtract 60s safety buffer)
	s.cachedToken = result.AccessToken
	s.tokenExpiry = time.Now().Add(14 * time.Minute)

	return result.AccessToken, nil
}

// transactionSignature generates the SNAP HMAC-SHA512 signature for transactions.
// Format: Method + ":" + EndpointPath + ":" + AccessToken + ":" + LowercaseHex(SHA256(Body)) + ":" + Timestamp
func (s *MidtransQRISService) transactionSignature(method, endpointPath, accessToken string, body []byte, ts string) string {
	bodyHash := sha256.Sum256(body)
	bodyHashHex := hex.EncodeToString(bodyHash[:])

	stringToSign := strings.ToUpper(method) + ":" + endpointPath + ":" +
		accessToken + ":" + bodyHashHex + ":" + ts

	mac := hmac.New(sha512.New, []byte(s.clientSecret))
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// ─── SNAP MPM QRIS Generation ────────────────────────────────────────────────

// CreateQRIS generates a QRIS barcode string and URL using SNAP QRIS MPM (/v1.0/qr/qr-mpm-generate).
func (s *MidtransQRISService) CreateQRIS(ctx context.Context, orderID string, amountRupiah int64) (*QRISResponse, error) {
	if orderID == "" {
		return nil, errors.New("orderID wajib diisi")
	}
	if len(orderID) > 36 {
		return nil, errors.New("orderID maksimal 36 karakter")
	}
	if amountRupiah <= 0 {
		return nil, errors.New("nominal harus lebih dari 0")
	}

	accessToken, err := s.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	externalID := orderID
	ts := s.timestamp()
	validUntil := time.Now().Add(15 * time.Minute).
		In(time.FixedZone("WIB", 7*60*60)).
		Format("2006-01-02T15:04:05-07:00")

	payload := QRISRequest{
		PartnerReferenceNo: externalID,
		Amount: Amount{
			Value:    fmt.Sprintf("%d.00", amountRupiah),
			Currency: "IDR",
		},
		MerchantID:     s.merchantId,
		ValidityPeriod: validUntil,
		AdditionalInfo: AdditionalInfo{
			Acquirer: "GOPAY",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	const endpointPath = "/v1.0/qr/qr-mpm-generate"
	signature := s.transactionSignature(http.MethodPost, endpointPath, accessToken, body, ts)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL()+endpointPath, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-TIMESTAMP", ts)
	req.Header.Set("X-SIGNATURE", signature)
	req.Header.Set("X-PARTNER-ID", s.partnerId)
	req.Header.Set("X-EXTERNAL-ID", externalID)
	req.Header.Set("CHANNEL-ID", s.channelId)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request QRIS: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Midtrans QRIS HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var result QRISResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode QRIS response: %w; body: %s", err, string(respBody))
	}

	if result.QRContent == "" && result.QRImage == "" && result.QRURL == "" {
		return nil, fmt.Errorf("response tidak berisi data QRIS: %s", string(respBody))
	}

	return &result, nil
}

// ─── domain.PaymentService Implementation ────────────────────────────────────

// GenerateQRCode implements domain.PaymentService.
// If SNAP credentials (clientId & clientSecret) are provided, it uses SNAP MPM QRIS.
// Otherwise, it falls back to Midtrans Core API /v2/charge.
func (s *MidtransQRISService) GenerateQRCode(req domain.QRCodePaymentRequest) (*domain.QRCodePaymentResponse, error) {
	// Mode 1: SNAP MPM QRIS (ASPI/BI Standard)
	if s.clientId != "" && s.clientSecret != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer cancel()

		qrisRes, err := s.CreateQRIS(ctx, req.OrderID, int64(req.Amount))
		if err != nil {
			return nil, err
		}

		qrImageURL := qrisRes.QRURL
		if qrImageURL == "" {
			qrImageURL = qrisRes.QRImage
		}

		validUntil := time.Now().Add(15 * time.Minute).
			In(time.FixedZone("WIB", 7*60*60)).
			Format("2006-01-02 15:04:05")

		return &domain.QRCodePaymentResponse{
			OrderID:       req.OrderID,
			Gateway:       domain.PaymentGatewayMidtrans,
			QRCodeURL:     qrImageURL,
			QRCodeString:  qrisRes.QRContent,
			ExpiresAt:     validUntil,
			Amount:        req.Amount,
			PaymentMethod: domain.PaymentMethodQRIS,
		}, nil
	}

	// Mode 2: Legacy Midtrans Core API /v2/charge (Fallback when SNAP keys are absent)
	payload := map[string]interface{}{
		"payment_type": "qris",
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": int64(req.Amount),
		},
		"qris": map[string]interface{}{
			"acquirer": "gopay",
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

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("midtrans request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("midtrans error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		StatusCode    string `json:"status_code"`
		StatusMessage string `json:"status_message"`
		TransactionID string `json:"transaction_id"`
		OrderID       string `json:"order_id"`
		QRString      string `json:"qr_string"`
		ExpiryTime    string `json:"expiry_time"`
		Actions       []struct {
			Name   string `json:"name"`
			Method string `json:"method"`
			URL    string `json:"url"`
		} `json:"actions"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	qrString := result.QRString
	qrImageURL := ""
	for _, action := range result.Actions {
		if action.Name == "generate-qr-code" {
			qrImageURL = action.URL
			break
		}
	}

	if qrString == "" && qrImageURL == "" {
		return nil, errors.New("midtrans did not return a QR code in response")
	}

	return &domain.QRCodePaymentResponse{
		OrderID:       req.OrderID,
		Gateway:       domain.PaymentGatewayMidtrans,
		QRCodeURL:     qrImageURL,
		QRCodeString:  qrString,
		ExpiresAt:     result.ExpiryTime,
		Amount:        req.Amount,
		PaymentMethod: domain.PaymentMethodQRIS,
	}, nil
}

// CheckPaymentStatus queries the payment status by order ID.
// Tries SNAP MPM status query if configured, and falls back to Core API status query.
func (s *MidtransQRISService) CheckPaymentStatus(orderID string) (string, error) {
	// If ServerKey is present, standard status endpoint /v2/{order_id}/status is the fastest and most reliable
	if s.serverKey != "" {
		url := fmt.Sprintf("%s/v2/%s/status", s.baseURL(), orderID)
		httpReq, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return "", err
		}
		httpReq.Header.Set("Accept", "application/json")
		httpReq.Header.Set("Authorization", "Basic "+s.basicAuth())

		resp, err := s.client.Do(httpReq)
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

	// SNAP MPM Query (/v1.0/qr/qr-mpm-query)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	accessToken, err := s.getAccessToken(ctx)
	if err != nil {
		return "", err
	}

	ts := s.timestamp()
	payload := QRISStatusQueryRequest{
		OriginalPartnerReferenceNo: orderID,
		ServiceCode:                "47",
		AdditionalInfo: AdditionalInfo{
			MerchantID: s.merchantId,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	const endpointPath = "/v1.0/qr/qr-mpm-query"
	signature := s.transactionSignature(http.MethodPost, endpointPath, accessToken, body, ts)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL()+endpointPath, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-TIMESTAMP", ts)
	req.Header.Set("X-SIGNATURE", signature)
	req.Header.Set("X-PARTNER-ID", s.partnerId)
	req.Header.Set("X-EXTERNAL-ID", orderID)
	req.Header.Set("CHANNEL-ID", s.channelId)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request QRIS query: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result QRISStatusQueryResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("decode QRIS query response: %w", err)
	}

	switch result.LatestTransactionStatus {
	case "00", "SUCCESS", "PAID":
		return "settlement", nil
	case "03", "PENDING":
		return "pending", nil
	case "04", "EXPIRED":
		return "expire", nil
	default:
		if result.TransactionStatusDesc != "" {
			return strings.ToLower(result.TransactionStatusDesc), nil
		}
		return result.LatestTransactionStatus, nil
	}
}
