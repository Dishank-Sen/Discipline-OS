package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTPInt() (int, error) {
	// 6-digit OTP: 100000 - 999999
	min := 100000
	max := 999999
	diff := max - min + 1

	n, err := rand.Int(rand.Reader, big.NewInt(int64(diff)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()) + min, nil
}