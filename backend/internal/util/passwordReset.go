package util

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTP() string {
	const otpLength = 6
	const digits = "0123456789"
	result := make([]byte, otpLength)
	charsetLength := big.NewInt(int64(len(digits)))

	for i := range result {
		num, _ := rand.Int(rand.Reader, charsetLength)
		result[i] = digits[num.Int64()]
	}

	return string(result)
}
