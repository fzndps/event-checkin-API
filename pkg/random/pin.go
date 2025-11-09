package random

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GeneratePIN mengasilkan random 4 digit PIN
func GeneratePIN() (string, error) {
	// Generate angka acak antara 0-9999
	max := big.NewInt(10000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("failed generate random number: %v", err)
	}

	// Format leading zeros
	pin := fmt.Sprintf("%04d", n.Int64())
	return pin, nil
}

// GenerateToken menghasilkan random token untuk qr code
// Format: 32 karakter hexa
func GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	token := fmt.Sprintf("%x", bytes)
	return token, nil
}
