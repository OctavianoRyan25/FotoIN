package infra

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"fotoin/backend/internal/config"
	"fotoin/backend/internal/domain"
)

const ProductDetail string = "Pembayaran Fotoyu"

type DuidkuSercvice struct {
	isProduction  bool
	apiKey        string
	merchantCode  string
	paymentMethod string
	callbackURL   string
	returnURL     string
	httpClient    *http.Client
}

type ItemDetail struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type DuidkuRequest struct {
	MerchantCode    string       `json:"merchantCode"`
	PaymentAmount   int          `json:"paymentAmount"`
	MerchantOrderId string       `json:"merchantOrderId"`
	ProductDetails  string       `json:"productDetails"`
	Email           string       `json:"email"`
	PaymentMethod   string       `json:"paymentMethod"`
	CustomerVaName  string       `json:"customerVaName"`
	ItemDetails     []ItemDetail `json:"itemDetails,omitempty"`
	ReturnURL       string       `json:"returnUrl"`
	CallbackURL     string       `json:"callbackUrl"`
	Signature       string       `json:"signature"`
	ExpiryPeriod    int          `json:"expiryPeriod,omitempty"`
}

type DuidkuResponse struct {
	MerchantCode  string `json:"merchantCode"`
	Reference     string `json:"reference"`
	PaymentUrl    string `json:"paymentUrl"`
	VaNumber      string `json:"vaNumber"`
	QrString      string `json:"qrString"`
	Amount        string `json:"amount"`
	StatusCode    string `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
}

type duidkuErrorResponse struct {
	Message string `json:"Message"`
}

type CreateOrderParams struct {
	MerchantOrderId string
	PaymentAmount   int
	Email           string
	CustomerVaName  string
	ItemDetails     []ItemDetail
	ReturnURL       string
	CallbackURL     string
	ExpiryPeriod    int
}

func NewDuidkuService(cfg config.DuitkuConfig) *DuidkuSercvice {
	return &DuidkuSercvice{
		isProduction:  cfg.Environment == "production",
		apiKey:        cfg.APIKey,
		merchantCode:  cfg.MerchantCode,
		paymentMethod: cfg.PaymentMethod,
		callbackURL:   cfg.CallbackURL,
		returnURL:     cfg.ReturnURL,
		httpClient:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *DuidkuSercvice) baseURL() string {
	if s.isProduction {
		return "https://passport.duitku.com/webapi/api/merchant/v2/inquiry"
	}
	return "https://sandbox.duitku.com/webapi/api/merchant/v2/inquiry"
}

func (s *DuidkuSercvice) CreateOrder(p CreateOrderParams) (*DuidkuResponse, error) {
	if p.MerchantOrderId == "" {
		return nil, fmt.Errorf("merchantOrderId wajib diisi")
	}
	if p.PaymentAmount <= 0 {
		return nil, fmt.Errorf("paymentAmount harus lebih besar dari 0")
	}
	if p.Email == "" {
		return nil, fmt.Errorf("email wajib diisi")
	}
	if p.CustomerVaName == "" {
		return nil, fmt.Errorf("customerVaName wajib diisi")
	}

	if len(p.ItemDetails) > 0 {
		total := 0
		for _, item := range p.ItemDetails {
			total += item.Price * item.Quantity
		}
		if total != p.PaymentAmount {
			return nil, fmt.Errorf("total itemDetails (%d) tidak sama dengan paymentAmount (%d)", total, p.PaymentAmount)
		}
	}

	sig, err := s.signature(p.MerchantOrderId, p.PaymentAmount)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat signature: %w", err)
	}

	reqBody := &DuidkuRequest{
		MerchantCode:    s.merchantCode,
		PaymentAmount:   p.PaymentAmount,
		MerchantOrderId: p.MerchantOrderId,
		ProductDetails:  ProductDetail,
		Email:           p.Email,
		PaymentMethod:   s.paymentMethod,
		CustomerVaName:  p.CustomerVaName,
		ItemDetails:     p.ItemDetails,
		ReturnURL:       p.ReturnURL,
		CallbackURL:     p.CallbackURL,
		Signature:       sig,
		ExpiryPeriod:    p.ExpiryPeriod,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("gagal marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, s.baseURL(), bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("gagal memanggil duitku: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response duitku: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp duidkuErrorResponse
		_ = json.Unmarshal(respBody, &errResp)
		return nil, fmt.Errorf("duitku merespon error %d: %s", resp.StatusCode, errResp.Message)
	}

	var result DuidkuResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("gagal parsing response duitku: %w", err)
	}

	return &result, nil
}

func (s *DuidkuSercvice) signature(merchantOrderId string, paymentAmount int) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("DUIDKU_APIKEY belum di-set")
	}
	stringToSign := fmt.Sprintf("%s%s%d", s.merchantCode, merchantOrderId, paymentAmount)

	mac := hmac.New(sha256.New, []byte(s.apiKey))
	if _, err := mac.Write([]byte(stringToSign)); err != nil {
		return "", err
	}

	return hex.EncodeToString(mac.Sum(nil)), nil
}

func (s *DuidkuSercvice) GenerateQRCode(req domain.QRCodePaymentRequest) (*domain.QRCodePaymentResponse, error) {
	if req.OrderID == "" {
		return nil, fmt.Errorf("order_id wajib diisi")
	}
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount harus lebih besar dari 0")
	}
	if req.CustomerEmail == "" {
		return nil, fmt.Errorf("customer_email wajib diisi")
	}

	result, err := s.CreateOrder(CreateOrderParams{
		MerchantOrderId: req.OrderID,
		PaymentAmount:   int(req.Amount),
		Email:           req.CustomerEmail,
		CustomerVaName:  req.CustomerEmail,
		ReturnURL:       strings.ReplaceAll(s.returnURL, "{order_id}", req.OrderID),
		CallbackURL:     s.callbackURL,
	})
	if err != nil {
		return nil, err
	}

	return &domain.QRCodePaymentResponse{
		OrderID:       req.OrderID,
		Gateway:       domain.PaymentGatewayDuidku,
		QRCodeString:  result.QrString,
		PaymentURL:    result.PaymentUrl,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
	}, nil
}

func (s *DuidkuSercvice) CheckPaymentStatus(orderID string) (string, error) {
	_ = orderID
	return "pending", nil
}

func (s *DuidkuSercvice) WebhookDuidku(p CreateOrderParams) error {
	return nil
}
