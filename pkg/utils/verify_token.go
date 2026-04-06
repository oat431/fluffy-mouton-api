package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateVerifyToken() string {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(tokenBytes)
}
