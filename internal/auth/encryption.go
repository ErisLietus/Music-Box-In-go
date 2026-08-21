package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HashEmail produces a consistent, deterministic lookup hash using a secret key.
func HashEmail(plainEmail, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(plainEmail))
	return hex.EncodeToString(h.Sum(nil))
}
