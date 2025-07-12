package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashSha256(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
