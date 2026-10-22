package infra

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"fotoin/backend/internal/config"
)

func generateTestRSAKey(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate rsa key: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal pkcs8: %v", err)
	}
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	}
	return string(pem.EncodeToMemory(block))
}

func TestMidtransQRISService_LoadRSAPrivateKey(t *testing.T) {
	pemStr := generateTestRSAKey(t)

	svc := NewMidtransQRISService(config.MidtransConfig{
		PrivateKey: pemStr,
	})

	rsaKey, err := svc.loadRSAPrivateKey()
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if rsaKey == nil {
		t.Fatal("expected rsaKey to be non-nil")
	}
}

func TestMidtransQRISService_TransactionSignature(t *testing.T) {
	svc := NewMidtransQRISService(config.MidtransConfig{
		ClientSecret: "secret123",
	})

	body := []byte(`{"partnerReferenceNo":"123"}`)
	sig := svc.transactionSignature("POST", "/v1.0/qr/qr-mpm-generate", "test-token", body, "2026-10-22T15:04:05+07:00")
	if sig == "" {
		t.Fatal("expected non-empty signature")
	}
}

func TestMidtransQRISService_Timestamp(t *testing.T) {
	svc := NewMidtransQRISService(config.MidtransConfig{})
	ts := svc.timestamp()
	if ts == "" {
		t.Fatal("expected non-empty timestamp")
	}
	_, err := time.Parse("2006-01-02T15:04:05-07:00", ts)
	if err != nil {
		t.Fatalf("timestamp %s is not valid ISO-8601: %v", ts, err)
	}
}
