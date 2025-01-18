package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateKey(length int) string {
	buffer := make([]byte, length)
	rand.Read(buffer)
	return base64.RawURLEncoding.EncodeToString(buffer)
}
