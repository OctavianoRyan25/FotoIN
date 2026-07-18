package handler

import "testing"

func TestIsDuidkuCallback(t *testing.T) {
	payload := map[string]any{
		"merchantOrderId": "order-123",
		"statusCode":      "00",
	}

	if !isDuidkuCallback(payload) {
		t.Fatal("expected payload to be recognized as Duitku callback")
	}
}

func TestIsSuccessfulDuidkuStatus(t *testing.T) {
	if !isSuccessfulDuidkuStatus("00") {
		t.Fatal("expected status code 00 to be treated as successful")
	}

	if isSuccessfulDuidkuStatus("01") {
		t.Fatal("expected non-success status to be rejected")
	}
}
