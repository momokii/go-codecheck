package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// The GenerateTokenSession function generates a random token of 32 bytes (256 bits) and returns it as
// a hexadecimal string.
func GenerateTokenSession() (string, error) {
	b := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
